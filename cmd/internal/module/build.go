package module

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"module-builder/internal/domain"
)

type PromoteType int

const (
	PromoteNone PromoteType = iota
	PromoteMinor
	PromoteMajor
)

func (t *PromoteType) Set(value string) error {
	switch value {
	case "none", "":
		*t = PromoteNone
	case "minor":
		*t = PromoteMinor
	case "major":
		*t = PromoteMajor
	default:
		return fmt.Errorf("only one of [<empty>, none, minor, major], given %s", value)
	}
	return nil
}

func (i PromoteType) String() string {
	switch i {
	case PromoteNone:
		return "None"
	case PromoteMinor:
		return "Minor"
	case PromoteMajor:
		return "Major"
	default:
		return ""
	}
}

type Config struct {
	LogWriter io.Writer   // logger
	Output    string      // where to put archives
	Dirs      []string    // module path (either abs or rel)
	Promote   PromoteType // type of promotion (dev, minor, major)
}

// Build archive and index for modules.
func Build(cfg Config) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("build recover: %v", err)
		}
	}()

	builder, err := newBuilder(cfg)
	if err != nil {
		_ = builder.Close() // sanity
		return err
	}

	defer builder.Close()
	return builder.Run()
}

type singleData struct {
	meta *os.File // metadata.yaml

	dir     string // module dir abs path
	dirBase string // module dir base

	hasChanges bool
}

type builder struct {
	logger *log.Logger

	archiveOutputDir string
	indexAbsPath     string

	modulesInfo []singleData

	promote PromoteType
}

func newBuilder(cfg Config) (*builder, error) {
	b := &builder{
		modulesInfo:      make([]singleData, len(cfg.Dirs)),
		promote:          cfg.Promote,
		logger:           log.New(cfg.LogWriter, "", log.Ltime|log.Lmicroseconds|log.Lshortfile),
		archiveOutputDir: cfg.Output,
	}

	// determine abs paths
	if err := b.collectAbsPaths(cfg.Dirs); err != nil {
		return nil, err
	}

	// determine if changes persist
	changes, err := b.getChanges(cfg.Dirs)
	if err != nil {
		return nil, err
	}

	// open files
	if err := b.openMetadataFiles(changes); err != nil {
		return b, err
	}

	return b, nil
}

func (b *builder) Run() error {
	modules := make([]domain.Module, len(b.modulesInfo))

	var merr error
	for i, m := range b.modulesInfo {
		tuple, err := b.bumpModuleMetaVersion(m)
		if err != nil {
			b.logger.Printf("ERROR: could not bump module %s version: %v", m.dirBase, err)
			merr = errors.Join(merr, err)
			continue
		}

		modules[i].NameVersionTuple = tuple
	}

	if merr != nil {
		b.logger.Printf("Error bumping modules versions: %v", merr)
		return fmt.Errorf("modules versions bump failed: %v", merr)
	}

	for i, module := range modules {
		shasum, err := makeArchive(module.Name, module.Version, b.archiveOutputDir)
		if err != nil {
			b.logger.Printf("ERROR: could make tgz with module %s: %v", module.NameVersionTuple, err)
			merr = errors.Join(merr, err)
			continue
		}

		modules[i].Sha256Sum = shasum
	}

	if merr != nil {
		b.logger.Printf("Error making tgz archives: %v", merr)
		return fmt.Errorf("archives baking failed: %v", merr)
	}

	b.logger.Printf("Updating index with %d modules", len(modules))
	if err := b.updateIndex(modules); err != nil {
		b.logger.Printf("Error updating index: %v", err)
		return fmt.Errorf("index update failed: %v", err)
	}

	return nil
}

func (b *builder) Close() error {
	if b == nil {
		return nil
	}

	var merr error
	for _, m := range b.modulesInfo {
		if err := m.meta.Close(); err != nil {
			merr = errors.Join(merr, err)
		}
	}

	return nil
}

func (b *builder) collectAbsPaths(dirs []string) error {
	ia, err := filepath.Abs(domain.IndexFileName)
	if err != nil {
		return fmt.Errorf("failed to determine abs path for the %s: %w", domain.IndexFileName, err)
	}
	b.indexAbsPath = ia

	if !filepath.IsAbs(b.archiveOutputDir) {
		archOutAbs, err := filepath.Abs(b.archiveOutputDir)
		if err != nil {
			return fmt.Errorf("failed to determine abs path for the %s: %w", b.archiveOutputDir, err)
		}

		b.archiveOutputDir = archOutAbs
	}

	for idx, dir := range dirs {
		absDir := dir
		if !filepath.IsAbs(dir) {
			var err error
			absDir, err = filepath.Abs(dir)
			if err != nil {
				return fmt.Errorf("failed to determine abs path for the %s: %w", dir, err)
			}
		}

		b.modulesInfo[idx] = singleData{
			dir:     absDir,
			dirBase: filepath.Base(absDir),
		}
	}

	return nil
}

func (b *builder) getChanges(dirs []string) ([]byte, error) {
	diffFlags := []string{
		"--exit-code",   // target the exit code
		"-z",            // simplier counting
		"--name-only",   // determine a change in a particular module
		"--no-ext-diff", // sanity
	}

	isChangeDetected := false

	cmd := exec.Command("git", "diff")
	cmd.Args = append(append(cmd.Args, diffFlags...), dirs...)

	output, err := cmd.Output()
	if err != nil {
		exitErr := new(exec.ExitError)
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			isChangeDetected = true
		} else {
			return nil, fmt.Errorf("failed execution of the '%s' command: %w", cmd.String(), err)
		}
	}

	// fail fast on incorrect cfg
	if isChangeDetected && b.promote != PromoteNone {
		return nil, fmt.Errorf("there are changes in modules, but promotion flag is provided")
	}

	return output, nil
}

func (b *builder) openMetadataFiles(data []byte) error {
	// parse modules names
	changedModules := map[string]struct{}{}
	for _, w := range parseModuleNames(data) {
		changedModules[string(w)] = struct{}{}
	}

	const metadataFileName = "metadata.yaml"
	for i, m := range b.modulesInfo {
		_, requiredChange := changedModules[m.dirBase]

		flags := os.O_RDONLY
		if requiredChange || b.promote != PromoteNone {
			flags = os.O_RDWR | os.O_CREATE
		}

		fileName := filepath.Join(m.dir, metadataFileName)
		f, err := os.OpenFile(fileName, flags, 0o644)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		m.meta = f
		m.hasChanges = requiredChange

		b.modulesInfo[i] = m
	}

	return nil
}

package sort

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"module-builder/internal/domain"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogWriter io.Writer // logger
}

// Index sorts index.yaml by name and version (if names are equal).
func Index(cfg Config) error {
	l := log.New(cfg.LogWriter, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)
	for _, indexFile := range [2]string{domain.ReleaseIndexFileName, domain.DevIndexFileName} {
		absIndexFile, err := filepath.Abs(indexFile)
		if err != nil {
			return fmt.Errorf("failed to determine abs path for the %s: %w", indexFile, err)
		}

		l.Printf("Sorting modules in the %s file", absIndexFile)
		if err := index(absIndexFile); err != nil {
			l.Printf("Error during sorting %s: %v", absIndexFile, err)
			return fmt.Errorf("sorting failed: %w", err)
		}
	}

	return nil
}

func index(name string) error {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	var index domain.HostOSConfigurationModules
	if err := yaml.NewDecoder(f).Decode(&index); err != nil {
		return fmt.Errorf("failed to deserialize data from %s: %w", name, err)
	}

	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate %s: %w", f.Name(), err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to set the offset in %s: %w", f.Name(), err)
	}

	slices.SortStableFunc(index.Spec.Modules, cmpModule)

	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	if err := enc.Encode(&index); err != nil {
		return fmt.Errorf("failed to serialize data to the %s: %w", name, err)
	}

	return enc.Close()
}

func cmpModule(a, b domain.Module) int {
	if a.Name == b.Name {
		return strings.Compare(a.Version, b.Version)
	}

	if a.Name < b.Name {
		return -1
	}

	return +1
}

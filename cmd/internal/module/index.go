package module

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"module-builder/internal/domain"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v3"
)

func (b *builder) updateDevIndex(newModules []domain.Module) error {
	indexFile, err := os.OpenFile(b.devIndexAbsPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("failed to read %s file: %w", b.devIndexAbsPath, err)
	}
	defer indexFile.Close()

	// create if did not exist
	if stat, _ := indexFile.Stat(); stat.Size() == 0 {
		if err := createIndex(indexFile, domain.DevHOCMObjName, newModules); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}

		return nil
	}

	var index domain.HostOSConfigurationModules
	if err := yaml.NewDecoder(indexFile).Decode(&index); err != nil {
		return fmt.Errorf("failed to deserialize %s: %w", b.devIndexAbsPath, err)
	}

	modulesBefore := slices.Clone(index.Spec.Modules)

	indexNameVer2Idx := make(map[string]int, len(index.Spec.Modules))
	for i, m := range index.Spec.Modules {
		indexNameVer2Idx[m.String()] = i
	}

	for _, newModule := range newModules {
		indexIdx, ok := indexNameVer2Idx[newModule.String()]
		if ok {
			if index.Spec.Modules[indexIdx].Sha256Sum == newModule.Sha256Sum {
				b.logger.Printf("Index is up to date for the module %s-%s, nothing to do", newModule.Name, newModule.Version)
				continue
			}

			// new module alredy exists, update to the new sha
			b.logger.Printf("Replacing hashsum for existing module %s-%s\n\t\tOld sha256: %s\n\t\tNew sha256: %s", newModule.Name, newModule.Version, newModule.Sha256Sum, newModule.Sha256Sum)

			index.Spec.Modules[indexIdx].Sha256Sum = newModule.Sha256Sum
			continue
		}

		b.logger.Printf("Appending module %s-%s, sha256: %s", newModule.Name, newModule.Version, newModule.Sha256Sum)
		index.Spec.Modules = append(index.Spec.Modules, newModule)
	}

	b.logger.Println("Trimming dev versions.")
	index.Spec.Modules = cutDevVersions(index.Spec.Modules, b.promote, newModules)

	// check changes
	if slices.EqualFunc(modulesBefore, index.Spec.Modules, func(a, b domain.Module) bool {
		return a.IsEqual(b)
	}) {
		b.logger.Println("Index has actual data, nothing to do.")
		return nil
	}

	if err := indexFile.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate %s: %w", indexFile.Name(), err)
	}
	if _, err := indexFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek %s: %w", indexFile.Name(), err)
	}

	enc := yaml.NewEncoder(indexFile)
	enc.SetIndent(2)
	if err := enc.Encode(&index); err != nil {
		return fmt.Errorf("failed to serialize data to the %s: %w", b.devIndexAbsPath, err)
	}

	return enc.Close()
}

func cutDevVersions(indexModules []domain.Module, promote PromoteType, newModules []domain.Module) []domain.Module {
	result := make([]domain.Module, 0, len(indexModules))
	newModulesNames := map[string]bool{}
	for _, m := range newModules {
		newModulesNames[m.Name] = true
	}

	if promote != PromoteNone {
		for _, m := range indexModules {
			if !newModulesNames[m.Name] || !strings.HasSuffix(m.Version, developmentTag) {
				result = append(result, m)
			}
		}

		return slices.Clip(result)
	}

	result = indexModules
	// drop all dev versions for same module which are less than our new modules versions
	for _, newModule := range newModules {
		for i, existingModule := range result {
			if existingModule.Name != newModule.Name ||
				!strings.HasSuffix(existingModule.Version, developmentTag) {
				continue
			}

			existingModuleVer, newModuleVer := semver.MustParse(existingModule.Version), semver.MustParse(newModule.Version)
			if existingModuleVer.LessThan(newModuleVer) {
				copy(result[i:], result[i+1:])
				result[len(result)-1] = domain.Module{}
				result = result[:len(result)-1]
			}
		}
	}

	return slices.Clip(result)
}

func createIndex(indexFile *os.File, indexName string, newModules []domain.Module) error {
	index := domain.HostOSConfigurationModules{
		APIVersion: "kaas.mirantis.com/v1alpha1",
		Kind:       "HostOSConfigurationModules",
		Metadata: struct {
			Name string "yaml:\"name\""
		}{
			Name: indexName,
		},
		Spec: struct {
			Modules []domain.Module "yaml:\"modules\""
		}{
			Modules: newModules,
		},
	}

	if err := indexFile.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate %s: %w", indexFile.Name(), err)
	}
	if _, err := indexFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek %s: %w", indexFile.Name(), err)
	}

	if err := yaml.NewEncoder(indexFile).Encode(&index); err != nil {
		return fmt.Errorf("failed to serialize index to %s: %w", indexFile.Name(), err)
	}

	return nil
}

func (b *builder) updateReleaseIndex() error {
	releaseIndexFile, err := os.OpenFile(b.releaseIndexAbsPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("failed to read %s file: %w", b.releaseIndexAbsPath, err)
	}
	defer releaseIndexFile.Close()

	devIndexFile, err := os.OpenFile(b.devIndexAbsPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("failed to read %s file: %w", b.devIndexAbsPath, err)
	}
	defer devIndexFile.Close()

	var devIndex domain.HostOSConfigurationModules
	if err := yaml.NewDecoder(devIndexFile).Decode(&devIndex); err != nil {
		return fmt.Errorf("failed to deserialize %s: %w", b.devIndexAbsPath, err)
	}

	var releaseModules []domain.Module
	for _, module := range devIndex.Spec.Modules {
		moduleVersion, err := semver.NewVersion(module.Version)
		if err != nil {
			return fmt.Errorf("failed to parse module version %s: %w", module.Version, err)
		}
		if moduleVersion.Prerelease() == "" {
			releaseModules = append(releaseModules, module)
		}
	}

	if err := createIndex(releaseIndexFile, domain.ReleaseHOCMObjName, releaseModules); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

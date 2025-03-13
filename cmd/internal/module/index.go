package module

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"module-builder/internal/domain"

	"gopkg.in/yaml.v3"
)

func (b *builder) updateDevIndex(newModules []domain.Module) error {
	indexFile, err := os.OpenFile(b.devIndexAbsPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("failed to read %s file: %w", b.devIndexAbsPath, err)
	}
	defer indexFile.Close()

	filteredModules := filterDevVersions(newModules)

	// create if did not exist
	if stat, _ := indexFile.Stat(); stat.Size() == 0 {
		if err := createIndex(indexFile, domain.DevHOCMObjName, filteredModules); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
		return nil
	}

	var index domain.HostOSConfigurationModules
	if err := yaml.NewDecoder(indexFile).Decode(&index); err != nil {
		return fmt.Errorf("failed to deserialize %s: %w", b.devIndexAbsPath, err)
	}

	index.Spec.Modules = filteredModules

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

func filterDevVersions(modules []domain.Module) []domain.Module {
	devModules := []domain.Module{}
	for _, newModule := range modules {
		if strings.HasSuffix(newModule.Version, developmentTag) {
			devModules = append(devModules, newModule)
		}
	}
	return slices.Clip(devModules)
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

func dropPromotedVersions(devIndexModules []domain.Module, newModules []domain.Module) []domain.Module {
	result := []domain.Module{}
	newModuleNames := make(map[string]struct{}, len(newModules))
	for _, newModule := range newModules {
		newModuleNames[newModule.Name] = struct{}{}
	}

	for _, oldModule := range devIndexModules {
		if _, exists := newModuleNames[oldModule.Name]; !exists {
			result = append(result, oldModule)
		}
	}
	return result
}

func (b *builder) promoteUpdateIndexes(newModules []domain.Module) error {
	releaseIndexFile, err := os.OpenFile(b.releaseIndexAbsPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("failed to read %s file: %w", b.releaseIndexAbsPath, err)
	}
	defer releaseIndexFile.Close()

	newProdModule := newModules
	if stat, _ := releaseIndexFile.Stat(); stat.Size() != 0 {
		var releaseIndex domain.HostOSConfigurationModules

		if err := yaml.NewDecoder(releaseIndexFile).Decode(&releaseIndex); err != nil {
			return fmt.Errorf("failed to deserialize %s: %w", b.releaseIndexAbsPath, err)
		}
		newProdModule = append(releaseIndex.Spec.Modules, newModules...)

		if err := releaseIndexFile.Truncate(0); err != nil {
			return fmt.Errorf("failed to truncate %s: %w", releaseIndexFile.Name(), err)
		}
		if _, err := releaseIndexFile.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek %s: %w", releaseIndexFile.Name(), err)
		}
	}
	if err := createIndex(releaseIndexFile, domain.ReleaseHOCMObjName, newProdModule); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	var devIndex domain.HostOSConfigurationModules
	devIndexFile, err := os.OpenFile(b.devIndexAbsPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("failed to read %s file: %w", b.devIndexAbsPath, err)
	}
	defer devIndexFile.Close()

	if err := yaml.NewDecoder(devIndexFile).Decode(&devIndex); err != nil {
		return fmt.Errorf("failed to deserialize %s: %w", b.devIndexAbsPath, err)
	}

	updatedDevModules := dropPromotedVersions(devIndex.Spec.Modules, newModules)

	if stat, _ := devIndexFile.Stat(); stat.Size() != 0 {
		if err := devIndexFile.Truncate(0); err != nil {
			return fmt.Errorf("failed to truncate %s: %w", devIndexFile.Name(), err)
		}
		if _, err := devIndexFile.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek %s: %w", devIndexFile.Name(), err)
		}
	}
	if err := createIndex(devIndexFile, domain.DevHOCMObjName, updatedDevModules); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	return nil
}

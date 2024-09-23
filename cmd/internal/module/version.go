package module

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"module-builder/internal/domain"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v3"
)

const developmentTag = "dev"

// bumpModuleMetaVersion parses metadata.yaml of a single module,
// and if required, bumps its version and modifies the metadata.yaml back.
func (b *builder) bumpModuleMetaVersion(data singleData) (meta domain.NameVersionTuple, _ error) {
	if err := yaml.NewDecoder(data.meta).Decode(&meta); err != nil {
		return meta, fmt.Errorf("failed to deserialize yaml %s: %w", data.meta.Name(), err)
	}

	moduleVersion, err := semver.NewVersion(meta.Version)
	if err != nil {
		b.logger.Printf("Malformed semver of the module %s: %v", meta, err)
		return meta, fmt.Errorf("failed to parse module version %s: %w", meta.Version, err)
	}

	var (
		mustModifyMeta       bool
		currentHasPrerelease = moduleVersion.Prerelease() != ""
	)

	if b.promote != PromoteNone && currentHasPrerelease {
		// if previous version is dev version, do promotion,
		// otherwise the current version is already promoted for release, so nothing to do.
		mustModifyMeta = true

		if b.promote == PromoteMajor {
			*moduleVersion = moduleVersion.IncMajor()
		} else {
			*moduleVersion = moduleVersion.IncMinor()
		}
	}

	if data.hasChanges { // promote is None
		mustModifyMeta = true

		*moduleVersion = moduleVersion.IncPatch()
		if currentHasPrerelease {
			*moduleVersion = moduleVersion.IncPatch() // incrementing on pre-version just drops it without increasing
		}

		*moduleVersion, err = moduleVersion.SetPrerelease(developmentTag)
		if err != nil {
			return meta, fmt.Errorf("failed to set prerelease version for module %s: %w", meta, err)
		}
	}

	if mustModifyMeta {
		newVersion := moduleVersion.String()
		if err := b.modifyMetadataVersion(data.meta, meta.Version, newVersion); err != nil {
			return meta, fmt.Errorf("modifying metadata.yaml failed: %w", err)
		}

		meta.Version = newVersion
	}

	return meta, nil
}

func (b *builder) modifyMetadataVersion(file *os.File, oldVersion, newVersion string) error {
	b.logger.Printf("Overwriting meta %s contents bumping version %s -> %s", file.Name(), oldVersion, newVersion)

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek %s: %w", file.Name(), err)
	}

	bb, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file.Name(), err)
	}

	var (
		contents      = bytes.Split(bb, []byte{'\n'})
		versionPrefix = []byte("version:")
	)
	for i, w := range contents {
		if bytes.HasPrefix(w, versionPrefix) {
			contents[i] = []byte("version: " + newVersion)
			break
		}
	}

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate %s: %w", file.Name(), err)
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek %s: %w", file.Name(), err)
	}

	if _, err := file.Write(bytes.Join(contents, []byte{'\n'})); err != nil {
		b.logger.Printf("Writing to file %s with new contents failed: %v", file.Name(), err)
		return fmt.Errorf("failed to write new contents to %s: %w", file.Name(), err)
	}

	return nil
}

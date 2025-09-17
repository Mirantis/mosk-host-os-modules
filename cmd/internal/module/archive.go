package module

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

// makeArchive makes a reproduceable tar-gzip archive with the module files and calculates its sha256sum.
func makeArchive(moduleName, moduleVersion, outputDir string) (string, error) {
	l := log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)

	tgzName := fmt.Sprintf("%s-%s.%s", filepath.Join(outputDir, moduleName), moduleVersion, "tgz")

	l.Printf("Starting to build the archive %s", tgzName)

	tmpFile, err := os.CreateTemp(outputDir, moduleName+"-"+moduleVersion+"-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	// lastly, move the file to the target destination
	defer func() {
		if err := os.Rename(tmpFile.Name(), tgzName); err != nil {
			l.Printf("failed to move %s to %s: %v", tmpFile.Name(), tgzName, err)
		}
	}()

	// firstly, close the tmp file
	defer func() {
		if err := tmpFile.Close(); err != nil {
			l.Printf("Error closing the temporary file %s: %v", tmpFile.Name(), err)
		}
	}()

	hash := sha256.New()

	mwr := io.MultiWriter(hash, tmpFile)

	if err := buildTarGz(moduleName, mwr); err != nil {
		l.Printf("Error building the archive %s: %v", tgzName, err)
		return "", fmt.Errorf("build the archive %s: %w", tgzName, err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func buildTarGz(root string, w io.Writer) error {
	gw := gzip.NewWriter(w)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	walkErr := filepath.Walk(root, func(filePath string, fileInfo fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Make path inside archive relative to root
		relPath, err := filepath.Rel(root, filePath)
		if err != nil {
			return fmt.Errorf("failed to compute relative path: %w", err)
		}
		// Normalize to forward slashes (POSIX style)
		relPath = filepath.ToSlash(relPath)

		// Special case for the root itself: skip adding
		if relPath == "." {
			return nil
		}

		// Create tar header from FileInfo
		header, err := tar.FileInfoHeader(fileInfo, "")
		if err != nil {
			return fmt.Errorf("failed to create header for %s: %w", filePath, err)
		}

		// Ensure reproducibility
		header.Name = relPath
		if fileInfo.IsDir() {
			// Directory entries should end with "/"
			header.Name += "/"
			header.Mode = 0o755
		} else {
			header.Mode = 0o600
		}
		header.ModTime = time.Unix(0, 0)
		header.ChangeTime = time.Unix(0, 0)
		header.AccessTime = time.Unix(0, 0)
		header.Uid = 0
		header.Gid = 0
		header.Gname = "root"
		header.Uname = "root"

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write header for %s: %w", relPath, err)
		}

		// Write file contents if it's a regular file
		if fileInfo.Mode().IsRegular() {
			file, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("failed to open %s: %w", filePath, err)
			}
			defer file.Close()

			if _, err := io.Copy(tw, file); err != nil {
				return fmt.Errorf("failed to copy file %s into archive: %w", relPath, err)
			}
		}

		return nil
	})

	if walkErr != nil {
		return fmt.Errorf("failed to build an archive: %w", walkErr)
	}

	return nil
}

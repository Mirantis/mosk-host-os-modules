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
	tw := tar.NewWriter(gw)

	// Iterate over files and add them to the tar archive
	walkErr := filepath.Walk(root, func(filePath string, fileInfo fs.FileInfo, err error) error {
		// Return on any error
		if err != nil {
			return err
		}
		// Skip directories and symlinks, pick only regular files
		if !fileInfo.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open %s to be archived: %w", filePath, err)
		}

		// Manage file header
		header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
		if err != nil {
			_ = file.Close()
			return fmt.Errorf("failed to create header for %s: %w", fileInfo.Name(), err)
		}

		// Make sure this tarball is reproducible by setting clean header for file
		header.ModTime = time.Unix(0, 0)
		header.ChangeTime = time.Unix(0, 0)
		header.AccessTime = time.Unix(0, 0)
		header.Mode = 0o600
		header.Uid = 0
		header.Gid = 0
		header.Gname = "root"
		header.Uname = "root"

		if err := tw.WriteHeader(header); err != nil {
			_ = file.Close()
			return fmt.Errorf("failed to write file header for %s: %w", header.Name, err)
		}

		// Copy file content to tar archive
		if _, err = io.Copy(tw, file); err != nil {
			return fmt.Errorf("failed to copy file %s into archive: %w", header.Name, err)
		}

		if err := file.Close(); err != nil {
			return fmt.Errorf("failed to close file %s: %w", header.Name, err)
		}

		return nil
	})

	if err := tw.Close(); err != nil {
		return fmt.Errorf("failed to close tar writer: %w", err)
	}
	if err := gw.Close(); err != nil {
		return fmt.Errorf("failed to close gunzip writer: %w", err)
	}

	if walkErr != nil {
		return fmt.Errorf("failed to build an archive: %w", walkErr)
	}

	return nil
}

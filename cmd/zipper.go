package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	yekazip "github.com/yeka/zip"
)

// ZipIfNeeded zips the file or directory if required (e.g., for password protection or directory upload)
func ZipIfNeeded(filePath string, password string) (string, error) {
	// Input validation
	if strings.TrimSpace(filePath) == "" {
		return "", fmt.Errorf("file path cannot be empty")
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("could not stat file %q: %w", filePath, err)
	}

	isDir := info.IsDir()
	shouldZip := isDir || password != ""

	if !shouldZip {
		// No zipping needed, just return the original file
		return filePath, nil
	}

	// Generate a better zip name
	zipName := generateUniqueZipName(filePath)

	// Check if zip file already exists
	if _, err := os.Stat(zipName); err == nil {
		return "", fmt.Errorf("zip file %q already exists", zipName)
	}

	zipFile, err := os.Create(zipName)
	if err != nil {
		return "", fmt.Errorf("could not create zip file %q: %w", zipName, err)
	}

	// Ensure proper cleanup with error handling
	var zipWriter *yekazip.Writer
	defer func() {
		if zipWriter != nil {
			if closeErr := zipWriter.Close(); closeErr != nil {
				// If we don't already have an error, use the close error
				if err == nil {
					err = fmt.Errorf("failed to close zip writer: %w", closeErr)
				}
			}
		}
		if closeErr := zipFile.Close(); closeErr != nil {
			// If we don't already have an error, use the close error
			if err == nil {
				err = fmt.Errorf("failed to close zip file: %w", closeErr)
			}
		}
		// Clean up zip file if there was an error
		if err != nil {
			os.Remove(zipName)
		}
	}()

	zipWriter = yekazip.NewWriter(zipFile)

	if isDir {
		// Walk the directory and add files
		err = filepath.Walk(filePath, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			// Skip directories and symlinks
			if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
				return nil
			}

			// Calculate relative path correctly
			relPath, relErr := filepath.Rel(filePath, path)
			if relErr != nil {
				return fmt.Errorf("could not calculate relative path for %q: %w", path, relErr)
			}

			return addFileToZip(zipWriter, path, relPath, password)
		})
		if err != nil {
			return "", fmt.Errorf("error zipping directory %q: %w", filePath, err)
		}
	} else {
		err = addFileToZip(zipWriter, filePath, filepath.Base(filePath), password)
		if err != nil {
			return "", fmt.Errorf("error zipping file %q: %w", filePath, err)
		}
	}

	return zipName, err
}

// addFileToZip adds a file to the zip writer, with optional password protection
func addFileToZip(zipWriter *yekazip.Writer, filePath, zipPath, password string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file %q: %w", filePath, err)
	}
	defer file.Close()

	var fw io.Writer
	if password != "" {
		fw, err = zipWriter.Encrypt(zipPath, password, yekazip.AES256Encryption)
		if err != nil {
			return fmt.Errorf("could not create encrypted entry for %q: %w", zipPath, err)
		}
	} else {
		fw, err = zipWriter.Create(zipPath)
		if err != nil {
			return fmt.Errorf("could not create zip entry for %q: %w", zipPath, err)
		}
	}

	// Use buffered copying for better performance with large files
	_, err = io.CopyBuffer(fw, file, make([]byte, 32*1024)) // 32KB buffer
	if err != nil {
		return fmt.Errorf("could not copy file %q to zip: %w", filePath, err)
	}

	return nil
}

// generateZipName creates a sensible zip file name from the input path
func generateZipName(filePath string) string {
	info, err := os.Stat(filePath)
	if err == nil && info.IsDir() {
		// If it's a directory, use the directory name
		return filePath + ".zip"
	}
	// Remove any existing extension for cleaner naming
	nameWithoutExt := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	return nameWithoutExt + ".zip"
}

// Alternative version that generates unique names if file exists
func generateUniqueZipName(filePath string) string {
	baseZipName := generateZipName(filePath)

	// If file doesn't exist, return the base name
	if _, err := os.Stat(baseZipName); os.IsNotExist(err) {
		return baseZipName
	}

	// Generate unique name with counter
	nameWithoutExt := strings.TrimSuffix(baseZipName, ".zip")
	counter := 1

	for {
		uniqueName := fmt.Sprintf("%s_%d.zip", nameWithoutExt, counter)
		if _, err := os.Stat(uniqueName); os.IsNotExist(err) {
			return uniqueName
		}
		counter++
	}
}

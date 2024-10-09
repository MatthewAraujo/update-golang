package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func extractTarGz(gzipPath, destDir string) error {
	log.Printf("Extracting %s to %s", gzipPath, destDir)

	file, err := os.Open(gzipPath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("error creating gzip reader: %v", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar file: %v", err)
		}

		targetPath := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("error creating directory: %v", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("error creating file: %v", err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("error writing file: %v", err)
			}

			// Set permissions for executables
			if filepath.Ext(targetPath) == "" {
				if err := os.Chmod(targetPath, 0755); err != nil {
					return fmt.Errorf("error setting executable permissions: %v", err)
				}
			}
		default:
			log.Printf("Skipping unknown type: %c in %s", header.Typeflag, targetPath)
		}
	}
	return nil
}

func removeOldGoFolder(dirPath string) error {
	log.Printf("Removing old Go folder: %s", dirPath)
	return os.RemoveAll(dirPath)
}

func downloadFile(url, filePath string) error {
	log.Printf("Downloading file from %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}

	return nil
}

func findTargetLine(content, target string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, target) {
			log.Printf("Target found in line %d: %s", i+1, line)
			return line
		}
	}
	return ""
}

func fetchPage(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error requesting Go page: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}
	return string(body)
}

func extractGoVersion(line string) string {
	fields := strings.Fields(line)
	if len(fields) > 0 {
		versionField := strings.Split(fields[2], `"`)
		if len(versionField) > 1 {
			log.Printf("Go version found: %s", versionField[1])
			return versionField[1]
		}
	}
	log.Fatal("Failed to extract Go version.")
	return ""
}

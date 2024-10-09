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

func fetchPage(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error requesting Go page:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading the response:", err)
	}
	return string(body)
}

func findTargetLine(content, target string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, target) {
			return line
		}
	}
	return ""
}

func extractGoVersion(text string) string {
	version := strings.Fields(text)
	if len(version) > 0 {
		version = strings.Split(version[2], `"`)
		return version[1]
	}
	log.Println("No version found in the provided text.")
	return ""
}

func downloadFile(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making GET request: %v", err)
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

func removeOldGoFolder(dirPath string) error {
	return os.RemoveAll(dirPath)
}

func extractTarGz(gzipPath, destDir string) error {
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
		}
	}

	return nil
}

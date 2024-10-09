package cmd

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Constants to facilitate configuration
const (
	GoDownloadPageURL   = "https://go.dev/dl/"
	DefaultGoInstallDir = "/usr/local/go"
)

// extractTarGz extracts a .tar.gz file to the specified destination directory.
func extractTarGz(gzipPath, destDir string) error {
	log.Printf("Extracting %s to %s", gzipPath, destDir)

	file, err := os.Open(gzipPath)
	if err != nil {
		return fmt.Errorf("error opening gzip file: %w", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("error creating gzip reader: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of tar file
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar file: %w", err)
		}

		targetPath := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("error creating directory %s: %w", targetPath, err)
			}
		case tar.TypeReg:
			// Ensure the directory exists
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("error creating directory for file %s: %w", targetPath, err)
			}

			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("error creating file %s: %w", targetPath, err)
			}

			// Copy data
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("error writing to file %s: %w", targetPath, err)
			}

			// Close the file after copying
			if err := outFile.Close(); err != nil {
				return fmt.Errorf("error closing file %s: %w", targetPath, err)
			}

			// Set permissions for executables
			if filepath.Ext(targetPath) == "" {
				if err := os.Chmod(targetPath, 0755); err != nil {
					return fmt.Errorf("error setting permissions for %s: %w", targetPath, err)
				}
			}
		default:
			log.Printf("Skipping unknown type: %c in %s", header.Typeflag, targetPath)
		}
	}

	return nil
}

// removeOldGoFolder removes the old Go installation.
func removeOldGoFolder(dirPath string) error {
	log.Printf("Removing old Go folder: %s", dirPath)
	err := os.RemoveAll(dirPath)
	if err != nil {
		return fmt.Errorf("error removing directory %s: %w", dirPath, err)
	}
	return nil
}

// downloadFile downloads a file from a URL to a local path.
func downloadFile(ctx context.Context, url, filePath string) error {
	log.Printf("Downloading file from %s", url)

	// Create a request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filePath, err)
	}
	defer func() {
		if cerr := out.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Copy data with a reasonable buffer size
	buf := make([]byte, 32*1024) // 32KB
	_, err = io.CopyBuffer(out, resp.Body, buf)
	if err != nil {
		return fmt.Errorf("error saving file %s: %w", filePath, err)
	}

	return nil
}

// findTargetLine finds the line that contains the target string in the content.
func findTargetLine(content, target string) (string, error) {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, target) {
			log.Printf("Target found in line %d: %s", i+1, line)
			return line, nil
		}
	}
	return "", errors.New("target not found in the provided lines")
}

// fetchPage performs an HTTP GET request to the provided URL and returns the response body as a string.
func fetchPage(ctx context.Context, url string) (string, error) {
	// Create a request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to access page: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading HTTP response: %w", err)
	}

	return string(body), nil
}

// extractGoVersion extracts the Go version from the provided line using regex.
func extractGoVersion(line string) (string, error) {
	// Regular expression to capture the Go version, e.g., go1.23.2
	re := regexp.MustCompile(`go\d+\.\d+\.\d+`)
	match := re.FindString(line)
	if match == "" {
		return "", errors.New("go version not found in the provided line")
	}
	log.Printf("Go version found: %s", match)
	return match, nil
}

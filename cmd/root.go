package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	versionFlag string
)

var rootCmd = &cobra.Command{
	Use:   "update-golang",
	Short: "A CLI tool to update your Go installation",
	Long:  `This CLI tool downloads the latest or specified version of Go, removes the old version, and installs the new one.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		var version string
		if versionFlag == "" {
			pageContent, err := fetchPage(ctx, GoDownloadPageURL)
			if err != nil {
				return fmt.Errorf("error fetching Go download page: %w", err)
			}

			line, err := findTargetLine(pageContent, `class="toggleVisible"`)
			if err != nil {
				return fmt.Errorf("error finding target line: %w", err)
			}

			version, err = extractGoVersion(line)
			if err != nil {
				return fmt.Errorf("error extracting Go version: %w", err)
			}
		} else {
			version = versionFlag
			log.Printf("Using specified version: %s", version)
		}

		log.Printf("Installing Go version: %s", version)

		downloadURL := fmt.Sprintf("https://golang.org/dl/%s.linux-amd64.tar.gz", version)
		filePath := fmt.Sprintf("%s.linux-amd64.tar.gz", version)

		// Download Go tarball if it doesn't exist
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if err := downloadFile(ctx, downloadURL, filePath); err != nil {
				return fmt.Errorf("error downloading file: %w", err)
			}
		} else {
			log.Printf("File %s already exists. Skipping download.", filePath)
		}

		if err := removeOldGoFolder(DefaultGoInstallDir); err != nil {
			return fmt.Errorf("error removing the old Go installation: %w", err)
		}

		if err := extractTarGz(filePath, "/usr/local"); err != nil {
			return fmt.Errorf("error extracting the file: %w", err)
		}

		if err := removeTarGzFile(filePath); err != nil {
			return fmt.Errorf("error removing the file: %w", err)
		}

		log.Println("Go installation successfully updated.")
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing the command: %v", err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&versionFlag, "version", "v", "", "Specify the Go version to install (e.g., go1.18). If not provided, the latest version will be used.")
}

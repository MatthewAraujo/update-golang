package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Root command
var rootCmd = &cobra.Command{
	Use:   "update-golang",
	Short: "A CLI tool to update your Go installation",
	Long:  `This CLI tool downloads the latest or specified version of Go, removes the old version, and installs the new version.`,
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetString("version")
		if version == "" {
			pageContent := fetchPage("https://go.dev/dl/")
			line := findTargetLine(pageContent, `class="toggleVisible"`)
			if line == "" {
				log.Fatal("Target not found.")
			}
			version = extractGoVersion(line)
		}

		log.Printf("Installing Go version: %s", version)

		downloadURL := fmt.Sprintf("https://go.dev/dl/%v.linux-amd64.tar.gz", version)
		filePath := version + ".linux-amd64.tar.gz"

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if err := downloadFile(downloadURL, filePath); err != nil {
				log.Fatalf("Error downloading the file: %v", err)
			}
		}

		goInstallDir := "/usr/local/go"
		if err := removeOldGoFolder(goInstallDir); err != nil {
			log.Fatalf("Error removing old Go folder: %v", err)
		}

		if err := extractTarGz(filePath, "/usr/local"); err != nil {
			log.Fatalf("Error extracting file: %v", err)
		}

		log.Println("Go installation updated successfully.")
	},
}

// Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Adding a flag for Go version
	rootCmd.Flags().StringP("version", "v", "", "Specify the Go version to install (e.g., go1.18). If not provided, the latest version will be used.")
}

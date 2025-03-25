package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// extractFileID extracts the Google Drive file ID from a raw share link.
// It supports the following formats:
//   - https://drive.google.com/file/d/FILE_ID/view?usp=sharing
//   - https://drive.google.com/open?id=FILE_ID
func extractFileID(rawLink string) (string, error) {
	// Attempt to extract using the "/file/d/<ID>" pattern.
	re := regexp.MustCompile(`file/d/([^/]+)`)
	matches := re.FindStringSubmatch(rawLink)
	if len(matches) > 1 {
		return matches[1], nil
	}

	// Fallback: parse the URL's query parameters to get "id".
	parsedURL, err := url.Parse(rawLink)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}
	if id := parsedURL.Query().Get("id"); id != "" {
		return id, nil
	}

	return "", fmt.Errorf("unable to extract file ID from URL: %s", rawLink)
}

func main() {
	fileName := "to_download.txt"
	maxParallel := 30

	// Create a context for the API calls.
	ctx := context.Background()

	// Initialize the Drive service using your credentials file.
	// Ensure you have a valid "credentials.json" in your working directory.
	driveService, err := drive.NewService(ctx, option.WithCredentialsFile("credentials.json"))
	if err != nil {
		log.Fatalf("Unable to create Drive service: %v", err)
	}

	// Open the text file containing raw links and destination paths.
	// Each file to download is represented by two consecutive lines:
	//   - First line: raw Google Drive link.
	//   - Second line: destination path.
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxParallel)

	// Process each pair of lines.
	for scanner.Scan() {
		rawLink := scanner.Text()
		if !scanner.Scan() { // Ensure there is a corresponding destination line.
			break
		}
		dest := scanner.Text()

		// Skip download if the destination file already exists.
		if _, err := os.Stat(dest); err == nil {
			fmt.Printf("File %s skipped because it exists already.\n", dest)
			continue
		}

		// Create parent directory if it doesn't exist.
		parentPath := filepath.Dir(dest)
		if _, err := os.Stat(parentPath); os.IsNotExist(err) {
			os.MkdirAll(parentPath, 0755)
		}

		// Enforce parallelism.
		sem <- struct{}{}
		wg.Add(1)

		go func(rawLink, dest string) {
			defer wg.Done()

			// Extract the file ID from the raw link.
			fileID, err := extractFileID(rawLink)
			if err != nil {
				log.Printf("Error extracting file ID from link %s: %v", rawLink, err)
				<-sem
				return
			}

			fmt.Printf("Downloading file ID %s to destination: %s\n", fileID, dest)
			resp, err := driveService.Files.Get(fileID).Download()
			if err != nil {
				log.Printf("Error downloading file %s: %v", fileID, err)
				<-sem
				return
			}
			defer resp.Body.Close()

			// Create the destination file.
			out, err := os.Create(dest)
			if err != nil {
				log.Printf("Error creating file %s: %v", dest, err)
				<-sem
				return
			}
			defer out.Close()

			// Stream the file download to the destination.
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				log.Printf("Error saving file %s: %v", dest, err)
			}
			<-sem
		}(rawLink, dest)
	}

	wg.Wait()
}

// HOW TO GET THE credentials.json:
// 1.	Go to Google Cloud Console:  https://console.cloud.google.com/welcome?inv=1&invt=Abs9mQ&project=golang-downloader
// Visit Google Cloud Console and sign in with your Google account.
// 	2.	Create or Select a Project:
// Either create a new project or select an existing one using the project dropdown at the top of the page.
// 	3.	Enable the Google Drive API:
// 	•	Navigate to APIs & Services > Library in the sidebar.
// 	•	Search for “Google Drive API”.
// 	•	Click on it and then click Enable.
// 	4.	Create Credentials:
// 	•	Navigate to APIs & Services > Credentials.
// 	•	Click Create Credentials and choose Service Account (for server-to-server use) or OAuth Client ID (if you need user consent).
// 	•	Follow the guided steps:
// 	•	For a Service Account:
// 	•	Provide a name and description.
// 	•	After creation, go to the service account’s Keys section, click Add Key > Create new key, choose JSON, and download the file.
// 	•	For an OAuth Client ID:
// 	•	You may need to configure the consent screen first.
// 	•	Select the appropriate application type (Desktop or Web) and download the credentials.
// 	5.	Save the Credentials File:
// Place the downloaded JSON file in your project’s working directory and name it credentials.json (or adjust the code to match your file name).

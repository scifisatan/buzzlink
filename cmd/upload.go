package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// UploadFile uploads the file at filePath to the buzzheavier API, with optional note.
// Returns the download link or an error.
func UploadFile(filePath, note string) (string, error) {
	filename := filepath.Base(filePath)
	uploadURL := "https://w.buzzheavier.com/" + filename
	if note != "" {
		encodedNote := base64.StdEncoding.EncodeToString([]byte(note))
		uploadURL += "?note=" + url.QueryEscape(encodedNote)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file for upload: %w", err)
	}
	defer file.Close()

	// Get file size for Content-Length
	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("could not stat file: %w", err)
	}
	fileSize := stat.Size()

	// Try to detect content type
	contentType := "application/octet-stream"
	if ext := filepath.Ext(filename); ext != "" {
		if mt := mime.TypeByExtension(ext); mt != "" {
			contentType = mt
		}
	}

	req, err := http.NewRequest("PUT", uploadURL, file)
	if err != nil {
		return "", fmt.Errorf("could not create upload request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)
	req.ContentLength = fileSize

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 201 {
		return "", fmt.Errorf("upload failed with status: %s", resp.Status)
	}

	var respBody struct {
		Code int `json:"code"`
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	decErr := json.Unmarshal(respBytes, &respBody)
	if decErr != nil {
		return "", fmt.Errorf("could not parse response: %w", decErr)
	}
	if respBody.Data.ID == "" {
		return "", fmt.Errorf("no ID in response")
	}

	downloadPageLink := "https://buzzheavier.com/" + respBody.Data.ID + "/download"

	// Make a GET request to the download page to get the hx-redirect header
	downloadLink, err := MakeBrowserLikeRequest(downloadPageLink)
	if err != nil {
		return "", err
	}

	return downloadLink, nil
}

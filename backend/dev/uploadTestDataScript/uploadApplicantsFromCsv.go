//go:build ignore

package main

// uploads airtable csv for applicants to mongo DB for testing
// cd backend
// go run dev/uploadTestDataScript/uploadApplicantsFromCsv.go

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kennedynguyen1/swipe-rank/backend/types"
)

func UploadApplicants() {
	// Open CSV file
	file, err := os.Open("dev/uploadTestDataScript/Fall 23 Rush App Responses.csv")
	if err != nil {
		fmt.Printf("Error opening CSV file: %v\n", err)
		return
	}
	defer file.Close()

	// Create CSV reader with custom configuration
	reader := csv.NewReader(file)
	reader.LazyQuotes = true      // Allow lazy quotes
	reader.FieldsPerRecord = -1   // Allow variable number of fields
	
	// Skip header row
	_, err = reader.Read()
	if err != nil {
		fmt.Printf("Error reading header: %v\n", err)
		return
	}

	// Process each row
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading record: %v\n", err)
			continue
		}

		// Create form response
		formData := types.FormResponses{
			Timestamp: time.Now().Format(time.RFC3339),
			Responses: []types.Response{
				{Question: "firstName", Answer: record[0]},
				{Question: "lastName", Answer: record[1]},
				{Question: "year", Answer: record[8]},
				{Question: "major", Answer: record[10]},
			},
		}

		// Handle file uploads (headshot, resume, cover letter)
		fileUrls := map[string]string{
			"image":       record[13],
			"resume":      record[14],
			"coverLetter": record[15],
		}

		for fileType, url := range fileUrls {
			if url == "" {
				continue
			}

			// Download file from URL
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Error downloading %s: %v\n", fileType, err)
				continue
			}
			defer resp.Body.Close()

			// Read file data
			fileData, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading %s data: %v\n", fileType, err)
				continue
			}

			// Get filename from URL
			filename := filepath.Base(url)
			
			// Determine MIME type based on file extension
			mimeType := "application/octet-stream"
			if strings.HasSuffix(strings.ToLower(filename), ".pdf") {
				mimeType = "application/pdf"
			} else if strings.HasSuffix(strings.ToLower(filename), ".png") {
				mimeType = "image/png"
			} else if strings.HasSuffix(strings.ToLower(filename), ".jpg") || strings.HasSuffix(strings.ToLower(filename), ".jpeg") {
				mimeType = "image/jpeg"
			}

			// Create file response
			fileResponse := types.Response{
				Question: fileType,
				Answer: map[string]interface{}{
					"type":     "file",
					"filename": filename,
					"mimeType": mimeType,
					"fileId":   []string{url}, // Using URL as fileId
					"data":     base64.StdEncoding.EncodeToString(fileData),
				},
			}
			formData.Responses = append(formData.Responses, fileResponse)
		}

		// Convert to JSON
		jsonData, err := json.Marshal(formData)
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			continue
		}

		// Send POST request to endpoint
		resp, err := http.Post("http://localhost:8080/formResponseListener", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		// Check response
		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Error response from server: %s\n", string(body))
			continue
		}

		fmt.Printf("Successfully uploaded application for %s %s\n", record[0], record[1])
	}
}

func main() {
	UploadApplicants()
}

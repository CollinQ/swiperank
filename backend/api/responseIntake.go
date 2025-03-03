package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kennedynguyen1/swipe-rank/backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)


func HandleFormResponse(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the JSON body
		var formData types.FormResponses
		if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
			http.Error(w, "Error parsing request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Create an ApplicantResponse instance
		applicant := types.Applicant{
			ID:          primitive.NewObjectID(),
			Rating:      0,  
			RatingCount: 0,  
			Timestamp:   formData.Timestamp,
		}

		// Create GridFS bucket
		bucket, err := gridfs.NewBucket(client.Database("akpsi-ucsb"))
		if err != nil {
			http.Error(w, "Error creating GridFS bucket: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Process responses and handle file uploads
		for _, resp := range formData.Responses {
			switch resp.Question {
			case "firstName":
				if strVal, ok := resp.Answer.(string); ok {
					applicant.FirstName = strVal
				}
			case "lastName":
				if strVal, ok := resp.Answer.(string); ok {
					applicant.LastName = strVal
				}
			case "major":
				if strVal, ok := resp.Answer.(string); ok {
					applicant.Major = strVal
				}
			case "year":
				if strVal, ok := resp.Answer.(string); ok {
					applicant.Year = strVal
				}
			case "coverLetter", "resume", "image":
				if answer, ok := resp.Answer.(map[string]interface{}); ok && answer["type"] == "file" {
					// Get data as string
					dataStr, ok := answer["data"].(string)
					if !ok {
						http.Error(w, "Invalid file data format", http.StatusBadRequest)
						return
					}
					
					// Decode base64 file data
					fileData, err := base64.StdEncoding.DecodeString(dataStr)
					if err != nil {
						http.Error(w, "Error decoding file data: "+err.Error(), http.StatusBadRequest)
						return
					}

					// Create unique filename with timestamp
					timestamp := time.Now().Unix()
					uniqueFileName := fmt.Sprintf("%d_%s", timestamp, answer["filename"].(string))

					// Upload to GridFS
					fileID, err := uploadToGridFS(bucket, uniqueFileName, fileData)
					if err != nil {
						http.Error(w, "Error uploading file: "+err.Error(), http.StatusInternalServerError)
						return
					}

					// Create FileInfo structure
					fileInfo := &types.FileInfo{
						FileID:      fileID.Hex(),
						FileName:    answer["filename"].(string),
						MimeType:    answer["mimeType"].(string),
						DriveFileID: answer["fileId"].([]interface{})[0].(string),
						UniqueName:  uniqueFileName,
						UploadedAt:  time.Now(),
					}

					switch resp.Question {
					case "coverLetter":
						applicant.CoverLetter = fileInfo
					case "resume":
						applicant.Resume = fileInfo
					case "image":
						applicant.Image = fileInfo
					}
				}
			}
		}

		// Get the collection
		collection := client.Database("akpsi-ucsb").Collection("applicants")

		// Insert the document
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := collection.InsertOne(ctx, applicant)
		if err != nil {
			http.Error(w, "Error inserting document: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Form response received successfully",
			"id":     result.InsertedID,
		})
		
		// Pretty print the application details
		prettyJSON, _ := json.MarshalIndent(applicant, "", "    ")
		fmt.Printf("Application Received:\n%s\n", string(prettyJSON))
	}
}

func uploadToGridFS(bucket *gridfs.Bucket, filename string, data []byte) (primitive.ObjectID, error) {
	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("error opening upload stream: %v", err)
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(data)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("error writing to stream: %v", err)
	}

	return uploadStream.FileID.(primitive.ObjectID), nil
} 
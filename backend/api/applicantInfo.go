package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"encoding/base64"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type FormResponse struct {
	Timestamp string `json:"timestamp"`
	Responses []struct {
		Question string      `json:"question"`
		Answer   interface{} `json:"answer"`
	} `json:"responses"`
}

// First, define the file answer structure
type FileAnswer struct {
	Type     string   `json:"type"`
	Data     string   `json:"data"`
	FileId   []string `json:"fileId"`
	Filename string   `json:"filename"`
	MimeType string   `json:"mimeType"`
}

func HandleFormResponse(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the JSON body
		var formData FormResponse
		if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
			http.Error(w, "Error parsing request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Process responses and handle file uploads
		processedResponses := make([]bson.M, 0)
		
		bucket, err := gridfs.NewBucket(client.Database("akpsi-ucsb"))
		if err != nil {
			http.Error(w, "Error initializing GridFS: "+err.Error(), http.StatusInternalServerError)
			return
		}
		for _, resp := range formData.Responses {
			processedResponse := bson.M{
				"question": resp.Question,
			}

			if answer, ok := resp.Answer.(map[string]interface{}); ok && answer["type"] == "file" {
				fileData, err := base64.StdEncoding.DecodeString(answer["data"].(string))
				if err != nil {
					http.Error(w, "Error decoding file data: "+err.Error(), http.StatusBadRequest)
					return
				}

				timestamp := time.Now().Unix()
				uniqueFileName := fmt.Sprintf("%d_%s", timestamp, answer["filename"].(string))

				fileID, err := uploadToGridFS(bucket, uniqueFileName, fileData)
				if err != nil {
					http.Error(w, "Error uploading file: "+err.Error(), http.StatusInternalServerError)
					return
				}

				processedResponse["answer"] = bson.M{
					"fileId":       fileID,
					"driveFileId":  answer["fileId"],
					"fileName":     answer["filename"].(string),
					"uniqueName":   uniqueFileName,
					"mimeType":     answer["mimeType"].(string),
					"type":         "file",
					"uploadedAt":   time.Now(),
				}
			} else {
				// For non-file responses, store the value directly
				processedResponse["answer"] = resp.Answer
			}

			processedResponses = append(processedResponses, processedResponse)
		}

		// Create document to insert
		document := bson.M{
			"submission_timestamp": formData.Timestamp,
			"responses":           processedResponses,
		}

		// Get the collection
		collection := client.Database("akpsi-ucsb").Collection("applicants")

		// Insert the document
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := collection.InsertOne(ctx, document)
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
		prettyJSON, _ := json.MarshalIndent(document, "", "    ")
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
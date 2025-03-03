package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/kennedynguyen1/swipe-rank/backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func GetLeastRatedApplicants(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers first, before any other operations
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000" 
		}
		w.Header().Set("Access-Control-Allow-Origin", frontendURL)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only allow GET method
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get the collection
		collection := client.Database("akpsi-ucsb").Collection("applicants")

		// Set up the context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// sort by ascending rating count and get first two applicants with the least ratings
		opts := options.Find().
			SetSort(bson.D{{Key: "ratingCount", Value: 1}}).
			SetLimit(2)

		// Execute the query
		cursor, err := collection.Find(ctx, bson.D{}, opts)
		if err != nil {
			http.Error(w, "Error finding applicants: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		// Decode the results into Applicant type
		var applicants []types.Applicant
		if err = cursor.All(ctx, &applicants); err != nil {
			http.Error(w, "Error decoding applicants: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Create GridFS bucket
		bucket, err := gridfs.NewBucket(client.Database("akpsi-ucsb"))
		if err != nil {
			http.Error(w, "Error creating GridFS bucket: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// For each applicant, fetch their files (image, cover letter, and resume)
		for i, applicant := range applicants {
			// Helper function to fetch and encode file data
			fetchFile := func(fileInfo *types.FileInfo) *types.FileInfo {
				if fileInfo != nil {
					fileID, err := primitive.ObjectIDFromHex(fileInfo.FileID)
					if err != nil {
						return nil // Skip if invalid ID
					}

					var buf bytes.Buffer
					_, err = bucket.DownloadToStream(fileID, &buf)
					if err != nil {
						return nil // Skip if download fails
					}

					fileInfo.Data = base64.StdEncoding.EncodeToString(buf.Bytes())
				}
				return fileInfo
			}

			// Update the applicant's files
			applicants[i].Image = fetchFile(applicant.Image)
			applicants[i].CoverLetter = fetchFile(applicant.CoverLetter)
			applicants[i].Resume = fetchFile(applicant.Resume)
		}

		// Return the results
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"applicants": applicants,
			},
		}); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Successfully retrieved applicant information!")
	}
}

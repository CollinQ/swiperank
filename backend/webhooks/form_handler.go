package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FormResponse struct {
	Timestamp string `json:"timestamp"`
	Responses []struct {
		Question string      `json:"question"`
		Answer   interface{} `json:"answer"`
	} `json:"responses"`
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

		// Create document to insert
		document := bson.M{
			"timestamp":  formData.Timestamp,
			"responses":  formData.Responses,
			"created_at": time.Now(),
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
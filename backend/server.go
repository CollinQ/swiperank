// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/joho/godotenv"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func main() {
// 	// Load .env file
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Get MongoDB URI from environment
// 	uri := os.Getenv("MONGODB_URI")
// 	if uri == "" {
// 		log.Fatal("MONGODB_URI not set in .env file")
// 	}

// 	// Configure client options
// 	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
// 	opts := options.Client().
// 		ApplyURI(uri).
// 		SetServerAPIOptions(serverAPI).
// 		SetConnectTimeout(10 * time.Second)

// 	// Create MongoDB client
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, opts)
// 	if err != nil {
// 		log.Fatalf("Connection error: %v", err)
// 	}
// 	defer client.Disconnect(ctx)

// 	// Verify connection
// 	if err := client.Ping(ctx, nil); err != nil {
// 		log.Fatalf("Ping failed: %v", err)
// 	}
// 	http.HandleFunc("/count-applicants", countApplicantsHandler(client))
// 	http.HandleFunc("/background-check", aiBackgroundCheck())
// 	fmt.Println("Successfully connected to MongoDB!")
// 	http.ListenAndServe(":8080", nil)
// }

// func aiBackgroundCheck() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		name := r.URL.Query().Get("name")

// 		w.Header().Set("Content-Type", "application/json")

// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(map[string]string{"message": "Hello, " + name})
// 	}

// }

// func countApplicantsHandler(client *mongo.Client) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()

// 		collection := client.Database("akpsi-ucsb").Collection("applicants")
// 		count, err := collection.CountDocuments(ctx, bson.M{})
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Error counting documents: %v", err), http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]int64{"count": count})
// 	}
// }
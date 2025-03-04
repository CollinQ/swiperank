//go:build ignore

package main

// clears applicant and applicant related collections for reset on database
// go run scripts/clearCollectionsScript/clearCollection.go -collections "applicants,fs.chunks,fs.files"

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Parse command line arguments
	collectionsStr := flag.String("collections", "", "Comma-separated names of collections to clear")
	flag.Parse()

	if *collectionsStr == "" {
		log.Fatal("Please provide collection names using -collections flag (comma-separated)")
	}

	collections := strings.Split(*collectionsStr, ",")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get MongoDB URI from environment
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI not set in .env file")
	}

	// Configure client options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPI).
		SetConnectTimeout(10 * time.Second)

	// Create MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	defer client.Disconnect(ctx)

	// Clear each collection
	for _, collectionName := range collections {
		collectionName = strings.TrimSpace(collectionName) // Remove any whitespace
		if collectionName == "" {
			continue
		}

		collection := client.Database("akpsi-ucsb").Collection(collectionName)
		result, err := collection.DeleteMany(ctx, bson.M{})
		if err != nil {
			log.Printf("Error clearing collection '%s': %v", collectionName, err)
			continue
		}

		log.Printf("Successfully deleted %d documents from collection '%s'", result.DeletedCount, collectionName)
	}
}

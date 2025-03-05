package main

import (
	"context"
	"log"
	"time"

	"backend/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func updater() {
    // Connect to MongoDB
    db.ConnectMongoDB("mongodb+srv://john:8900@ucsb-akpsi.18g4e.mongodb.net/?retryWrites=true&w=majority&appName=UCSB-AKPsi")

    // Get the applicants collection
    collection := db.GetCollection("applicants")

    // Create a context with a timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Define the update operation
    update := bson.M{
        "$set": bson.M{
            "matches_played": []primitive.ObjectID{},
            "elo":            1500,
        },
    }

    // Perform the update operation on all documents in the collection
    result, err := collection.UpdateMany(ctx, bson.M{}, update)
    if err != nil {
        log.Fatal("Failed to update applicants: ", err)
    }

    log.Printf("Matched %d documents and updated %d documents.\n", result.MatchedCount, result.ModifiedCount)
}
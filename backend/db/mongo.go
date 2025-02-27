package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// making one connection to the mongo for all endpoints to use
var Client *mongo.Client

func ConnectMongoDB(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB not responding: ", err)
	}

	log.Println("Connected to MongoDB")
	Client = client
}

func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("akpsi-ucsb").Collection(collectionName)
}
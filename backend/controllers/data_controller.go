package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"backend/db"
	"backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DataController struct {
	collection *mongo.Collection
}

func NewDataController() *DataController {
	return &DataController{
		collection: db.GetCollection("data"),
	}
}

func (dc *DataController) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := dc.collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		log.Println("MongoDB Find data error: ", err)
		return
	}
	defer cursor.Close(ctx)

	var data []models.Data
	if err = cursor.All(ctx, &data); err != nil {
		http.Error(w, "Error decoding data", http.StatusInternalServerError)
		log.Println("Cursor decode error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

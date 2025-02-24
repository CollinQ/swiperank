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
)

type ProjectController struct {
	collection *mongo.Collection
}

func NewProjectController() *ProjectController {
	return &ProjectController{
		collection: db.GetCollection("projects"),
	}
}

func (pc *ProjectController) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := pc.collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		log.Println("MongoDB Find project error: ", err)
		return
	}
	defer cursor.Close(ctx)

	var projects []models.Project
	if err = cursor.All(ctx, &projects); err != nil {
		http.Error(w, "Error decoding projects", http.StatusInternalServerError)
		log.Println("Cursor decode error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
} 
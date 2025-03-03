package routes

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

func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	
	collection := db.GetCollection("projects")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find all documents in the projects collection (empty filter means all)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		log.Println("MongoDB Find project error: ", err)
		return
	}
	defer cursor.Close(ctx)

	var projects []models.Project
	err = cursor.All(ctx, &projects)
	if err != nil {
		http.Error(w, "Error decoding projects", http.StatusInternalServerError)
		log.Println("Cursor decode error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}
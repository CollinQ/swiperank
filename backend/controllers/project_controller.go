package controllers

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"

	"backend/db"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (pc *ProjectController) Create(w http.ResponseWriter, r *http.Request) {
	var project models.Project

	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}
	if project.Name == "" {
		http.Error(w, "Project name is required", http.StatusBadRequest)
		return
	}

	project.ID = primitive.NewObjectID()
	project.CompletedComparisons = 0

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	project.TotalComparisons = calculateSwissTotalComparisons(int(150)) // NEED APPLICATIONS IN MONGODB

	_, err := pc.collection.InsertOne(ctx, project)
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		log.Println("mongoDB Insert project error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

func calculateSwissTotalComparisons(numApplicants int) int {
	if numApplicants < 2 {
		return 0
	}

	rounds := int(math.Ceil(math.Log2(float64(numApplicants)))) + 2
	return rounds * (numApplicants / 2)
}
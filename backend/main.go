package main

import (
	"log"
	"net/http"
	"os"


	"backend/db"
	"backend/routes"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env failed to load")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("No MONGO_URI found in .env")
	}

	db.ConnectMongoDB(mongoURI)
	router := chi.NewRouter()
	router.Get("/api/projects", routes.GetProjectsHandler)

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", router)
}
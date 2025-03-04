package main

import (
	"log"
	"net/http"
	"os"

	"backend/db"
	"backend/routes"

	// "backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env failed to load")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("No MONGODB_URI found in .env")
	}

	// clerkAPIkey := os.Getenv("CLERK_API_KEY")
	// if clerkAPIkey == "" {
	// 	log.Fatal("No CLERK_API_KEY found in .env")
	// }

	// middleware.InitClerk(clerkAPIkey)

	db.ConnectMongoDB(mongoURI)

	router := chi.NewRouter()
	routes.SetupRoutes(router)
	
	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", router)
}
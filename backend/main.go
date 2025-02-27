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
	"github.com/go-chi/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
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

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	// router.With(middleware.AuthMiddleware).Get("/api/projects", routes.GetProjectsHandler)
	routes.SetupRoutes(router)
	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", router)
}
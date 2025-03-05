package routes

import (
	"encoding/json"
	"net/http"
	"os"

	"backend/controllers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRoutes(router chi.Router) {
	// Setup CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{func() string {
			url := os.Getenv("FRONTEND_URL")
			if url == "" {
				return "http://localhost:3000"
			}
			return url
		}()},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize controllers
	projectController := controllers.NewProjectController()
	applicantController := controllers.NewApplicantController()
	formResponseController := controllers.NewFormResponseController()
	// dataController := controllers.NewDataController()

	router.Route("/api", func(r chi.Router) {
		// Project routes
		r.Get("/projects", projectController.GetAll)
		// r.Get("/data", dataController.GetAll) // TODO // when clicking "ADD NEW PROJECT" I want this to display all new projects, NOT NECESSARY FOR NOW. FOCUS ON MAKING ONE WORK
		r.Post("/projects", projectController.Create)

		// r.Get("/applicants", applicantController.GetAll) // TODO
		r.Get("/applicants", applicantController.GetById)


		r.Get("/getTwoForComparison", applicantController.GetTwoForComparison)
		r.Post("/updateElo", applicantController.UpdateElo)
		r.Get("/rankings", applicantController.GetRankings)
		// Additional routes from server.go
		r.Get("/background-check", aiBackgroundCheck())
		r.Post("/formResponseListener", formResponseController.HandleFormResponse)
	})
}

func aiBackgroundCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello, " + name})
	}
}



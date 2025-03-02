package routes

import (
	"backend/controllers"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router chi.Router) {
	// Initialize controllers
	projectController := controllers.NewProjectController()
	applicantController := controllers.NewApplicantController()
	// dataController := controllers.NewDataController()

	router.Route("/api", func(r chi.Router) {
		// Project routes
		r.Get("/projects", projectController.GetAll)
		// r.Get("/data", dataController.GetAll) // TODO // when clicking "ADD NEW PROJECT" I want this to display all new projects, NOT NECESSARY FOR NOW. FOCUS ON MAKING ONE WORK
		r.Post("/projects", projectController.Create)

		r.Get("/applicants", applicantController.GetAll) // TODO
		r.Get("/applicants/{id}", applicantController.GetById)

		r.Post("/applicants", applicantController.Create) // Not too sure for now, my idea is that this will be used to add applicants to a project
		r.Put("/applicants/{id}", applicantController.Update) // this will be used to update an applicant's information, so what I want to do is update a counter for this applicant, and then you take whoever counter got updated to do the swiss comparison.

		r.Get("/getTwoForComparison", applicantController.GetTwoForComparison)
		r.Post("/updateElo", applicantController.UpdateElo)
		// Add other routes here
	})
} 
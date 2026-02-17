package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"courseshare/config"
	"courseshare/handlers"
	"courseshare/models"
)




func main() {

	config.ConnectDatabase()

	config.DB.AutoMigrate(&models.Course{}, &models.Note{})

r := chi.NewRouter()

r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Use(middleware.StripSlashes)



	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server is running"))
})


	// Course routes
	r.Get("/courses", handlers.GetCourses)
	r.Post("/courses", handlers.CreateCourse)

	// Note routes
	r.Get("/courses/{id}/notes", handlers.GetNotesByCourse)
	r.Get("/notes/{id}", handlers.GetNote)
	r.Post("/notes", handlers.CreateNote)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}

package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"courseshare/config"
	"courseshare/handlers"
	"courseshare/models"
)

func main() {

	// Connect Database
	config.ConnectDatabase()

	// Auto migrate tables
	config.DB.AutoMigrate(&models.Course{}, &models.Note{})

	// Router
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	// CORS Middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	}))

	// Routes
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
	http.ListenAndServe("0.0.0.0:8080", r)
}
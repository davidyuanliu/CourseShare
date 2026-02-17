package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/davidyuanliu/CourseShare/backend/config"
	"github.com/davidyuanliu/CourseShare/backend/models"
	"github.com/davidyuanliu/CourseShare/backend/handlers"
)

func main() {

	config.ConnectDatabase()

	config.DB.AutoMigrate(&models.Course{}, &models.Note{})

	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Incoming:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
})


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

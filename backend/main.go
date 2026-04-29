package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"courseshare/config"
	"courseshare/handlers"
	authMiddleware "courseshare/middleware"
	"courseshare/models"
)

func main() {

	// Connect Database
	config.ConnectDatabase()

	// Auto migrate tables
	config.DB.AutoMigrate(&models.User{}, &models.Course{}, &models.Note{}, &models.HelpfulVote{}, &models.SavedNote{})

	// Router
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	allowedOrigins := getAllowedOrigins()

	// CORS Middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	// Auth routes
	r.Post("/auth/register", handlers.RegisterUser)
	r.Post("/auth/login", handlers.LoginUser)

	// Course routes
	r.Get("/courses", handlers.GetCourses)
	r.Post("/courses", handlers.CreateCourse)

	// Note routes
	r.Get("/courses/{id}/notes", handlers.GetNotesByCourse)
	r.Get("/notes/{id}", handlers.GetNote)
	r.Get("/my-notes", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.ExtractUserFromToken(http.HandlerFunc(handlers.GetMyNotes)).ServeHTTP(w, r)
	})
	r.Get("/saved-notes", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.ExtractUserFromToken(http.HandlerFunc(handlers.GetSavedNotes)).ServeHTTP(w, r)
	})
	r.Post("/notes/{id}/helpful", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.ExtractUserFromToken(http.HandlerFunc(handlers.MarkNoteHelpful)).ServeHTTP(w, r)
	})
	r.Delete("/notes/{id}/helpful", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.ExtractUserFromToken(http.HandlerFunc(handlers.UnmarkNoteHelpful)).ServeHTTP(w, r)
	})
	r.Post("/notes/{id}/save", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.ExtractUserFromToken(http.HandlerFunc(handlers.SaveNote)).ServeHTTP(w, r)
	})
	r.Delete("/notes/{id}/save", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.ExtractUserFromToken(http.HandlerFunc(handlers.UnsaveNote)).ServeHTTP(w, r)
	})
	r.Post("/notes", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.ExtractUserFromToken(http.HandlerFunc(handlers.CreateNote)).ServeHTTP(w, r)
	})
	r.Put("/notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.ExtractUserFromToken(http.HandlerFunc(handlers.UpdateNote)).ServeHTTP(w, r)
	})
	r.Delete("/notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.ExtractUserFromToken(http.HandlerFunc(handlers.DeleteNote)).ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := "0.0.0.0:" + port
	log.Printf("Server running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func getAllowedOrigins() []string {
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	if originsEnv == "" {
		return []string{"http://localhost:4200", "http://localhost:3000"}
	}

	rawOrigins := strings.Split(originsEnv, ",")
	origins := make([]string, 0, len(rawOrigins))
	for _, origin := range rawOrigins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	if len(origins) == 0 {
		return []string{"http://localhost:4200", "http://localhost:3000"}
	}

	return origins
}

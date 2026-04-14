package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"courseshare/config"
	"courseshare/models"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var register struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&register); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate required fields
	missingFields := []string{}
	if strings.TrimSpace(register.Email) == "" {
		missingFields = append(missingFields, "email")
	}
	if strings.TrimSpace(register.Password) == "" {
		missingFields = append(missingFields, "password")
	}

	if len(missingFields) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":         "Missing required fields",
			"missingFields": missingFields,
		})
		return
	}

	// Validate email format
	email := strings.TrimSpace(register.Email)
	if _, err := mail.ParseAddress(email); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid email format",
		})
		return
	}

	// Check for duplicate email (case-insensitive)
	var existingUser models.User
	result := config.DB.Where("LOWER(email) = LOWER(?)", email).First(&existingUser)
	if result.Error == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Email already registered",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create user
	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&login); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate required fields
	missingFields := []string{}
	if strings.TrimSpace(login.Email) == "" {
		missingFields = append(missingFields, "email")
	}
	if strings.TrimSpace(login.Password) == "" {
		missingFields = append(missingFields, "password")
	}

	if len(missingFields) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":         "Missing required fields",
			"missingFields": missingFields,
		})
		return
	}

	// Find user by email (case-insensitive)
	var user models.User
	result := config.DB.Where("LOWER(email) = LOWER(?)", login.Email).First(&user)
	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid credentials",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid credentials",
		})
		return
	}

	// Generate JWT token
	tokenString, err := generateJWT(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return success response with token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func generateJWT(userID uint, email string) (string, error) {
	// Get secret from environment or use default for development
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-key-please-change-in-production"
	}

	// Create JWT claims
	claims := jwt.MapClaims{
		"sub":   userID, // Subject = user ID
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	// Create and sign token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

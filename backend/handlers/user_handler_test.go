package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"courseshare/config"
	"courseshare/models"
)

func TestRegisterUserSuccess(t *testing.T) {
	setupTestDB(t)

	body := strings.NewReader(`{"email":"newuser@example.com","password":"securepass123"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	RegisterUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var resp struct {
		Message string `json:"message"`
		User    struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "User registered successfully" {
		t.Fatalf("unexpected message: %s", resp.Message)
	}

	if resp.User.Email != "newuser@example.com" {
		t.Fatalf("expected email newuser@example.com, got %s", resp.User.Email)
	}

	// Verify user was created in database
	var user models.User
	result := config.DB.First(&user, resp.User.ID)
	if result.Error != nil {
		t.Fatalf("user not found in database: %v", result.Error)
	}
}

func TestRegisterUserMissingFields(t *testing.T) {
	setupTestDB(t)

	tests := []struct {
		name           string
		body           string
		expectedFields []string
	}{
		{
			name:           "missing both email and password",
			body:           `{"email":"","password":""}`,
			expectedFields: []string{"email", "password"},
		},
		{
			name:           "missing email",
			body:           `{"email":"","password":"validpass123"}`,
			expectedFields: []string{"email"},
		},
		{
			name:           "missing password",
			body:           `{"email":"test@example.com","password":""}`,
			expectedFields: []string{"password"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			RegisterUser(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d", rr.Code)
			}

			var resp map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			missing, ok := resp["missingFields"].([]interface{})
			if !ok {
				t.Fatalf("missingFields not returned")
			}

			if len(missing) != len(tt.expectedFields) {
				t.Fatalf("expected %d missing fields, got %d", len(tt.expectedFields), len(missing))
			}
		})
	}
}

func TestRegisterUserInvalidEmail(t *testing.T) {
	setupTestDB(t)

	invalidEmails := []string{
		"notanemail.com",
		"@example.com",
		"user@",
		"user name@example.com",
	}

	for _, invalidEmail := range invalidEmails {
		t.Run(fmt.Sprintf("email: %s", invalidEmail), func(t *testing.T) {
			body := strings.NewReader(fmt.Sprintf(`{"email":"%s","password":"validpass123"}`, invalidEmail))
			req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			RegisterUser(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d", rr.Code)
			}

			var resp map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp["error"] != "Invalid email format" {
				t.Fatalf("expected 'Invalid email format' error, got %s", resp["error"])
			}
		})
	}
}

func TestRegisterUserDuplicateEmail(t *testing.T) {
	setupTestDB(t)

	// Create first user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
	user := models.User{
		Email:    "existing@example.com",
		Password: string(hashedPassword),
	}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	// Try to register with the same email
	body := strings.NewReader(`{"email":"existing@example.com","password":"password2"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	RegisterUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "Email already registered" {
		t.Fatalf("expected 'Email already registered' error, got %s", resp["error"])
	}
}

func TestRegisterUserDuplicateEmailCaseInsensitive(t *testing.T) {
	setupTestDB(t)

	// Create first user with lowercase email
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
	user := models.User{
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	// Try to register with uppercase version of same email
	body := strings.NewReader(`{"email":"TEST@EXAMPLE.COM","password":"password2"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	RegisterUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "Email already registered" {
		t.Fatalf("expected 'Email already registered' error, got %s", resp["error"])
	}
}

func TestRegisterUserPasswordHashed(t *testing.T) {
	setupTestDB(t)

	plainPassword := "mySecurePassword123!"
	body := strings.NewReader(fmt.Sprintf(`{"email":"hashtest@example.com","password":"%s"}`, plainPassword))
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	RegisterUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var resp struct {
		User struct {
			ID uint `json:"id"`
		} `json:"user"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify password is hashed in database
	var user models.User
	if err := config.DB.First(&user, resp.User.ID).Error; err != nil {
		t.Fatalf("user not found: %v", err)
	}

	// Password should not match plaintext
	if user.Password == plainPassword {
		t.Fatalf("password is not hashed - stored as plaintext")
	}

	// Password hash should be valid with bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword)); err != nil {
		t.Fatalf("password hash is invalid: %v", err)
	}
}

func TestLoginUserSuccess(t *testing.T) {
	setupTestDB(t)

	// Create a user first
	email := "login@example.com"
	password := "testpassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Login with correct credentials
	body := strings.NewReader(fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password))
	req := httptest.NewRequest(http.MethodPost, "/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	LoginUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp struct {
		Message string `json:"message"`
		Token   string `json:"token"`
		User    struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "Login successful" {
		t.Fatalf("unexpected message: %s", resp.Message)
	}

	if resp.Token == "" {
		t.Fatalf("no token returned")
	}

	if resp.User.ID != user.ID {
		t.Fatalf("expected user id %d, got %d", user.ID, resp.User.ID)
	}
}

func TestLoginUserInvalidEmail(t *testing.T) {
	setupTestDB(t)

	body := strings.NewReader(`{"email":"nonexistent@example.com","password":"anypassword"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	LoginUser(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "Invalid credentials" {
		t.Fatalf("expected 'Invalid credentials' error, got %s", resp["error"])
	}
}

func TestLoginUserWrongPassword(t *testing.T) {
	setupTestDB(t)

	// Create a user
	email := "user@example.com"
	password := "correctpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Try to login with wrong password
	body := strings.NewReader(fmt.Sprintf(`{"email":"%s","password":"wrongpassword"}`, email))
	req := httptest.NewRequest(http.MethodPost, "/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	LoginUser(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "Invalid credentials" {
		t.Fatalf("expected 'Invalid credentials' error, got %s", resp["error"])
	}
}

func TestLoginUserMissingFields(t *testing.T) {
	setupTestDB(t)

	tests := []struct {
		name           string
		body           string
		expectedFields []string
	}{
		{
			name:           "missing both email and password",
			body:           `{"email":"","password":""}`,
			expectedFields: []string{"email", "password"},
		},
		{
			name:           "missing email",
			body:           `{"email":"","password":"validpass123"}`,
			expectedFields: []string{"email"},
		},
		{
			name:           "missing password",
			body:           `{"email":"test@example.com","password":""}`,
			expectedFields: []string{"password"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", body)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			LoginUser(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d", rr.Code)
			}

			var resp map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			missing, ok := resp["missingFields"].([]interface{})
			if !ok {
				t.Fatalf("missingFields not returned")
			}

			if len(missing) != len(tt.expectedFields) {
				t.Fatalf("expected %d missing fields, got %d", len(tt.expectedFields), len(missing))
			}
		})
	}
}

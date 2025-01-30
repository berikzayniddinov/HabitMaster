package handlers_test

import (
	"HabitMaster/handlers"
	"bytes"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testDB *sql.DB

func setupTestDB(t *testing.T) {
	var err error
	testDB, err = sql.Open("postgres", "user=habitmaster password=123456 dbname=dbhabit sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Очистка таблицы перед тестами
	_, err = testDB.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("Failed to clean up test database: %v", err)
	}
}

func teardownTestDB(t *testing.T) {
	testDB.Close()
}

func TestCreateUser(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := map[string]string{
		"name":     "Test User",
		"email":    "testuser@example.com",
		"password": "securepassword",
	}
	body, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := handlers.CreateUser(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", recorder.Code)
	}
}

func TestGetUsers(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	_, err := testDB.Exec(`INSERT INTO users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`, "Test User", "testuser_unique@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	req, err := http.NewRequest("GET", "/api/get-users", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	recorder := httptest.NewRecorder()
	handler := handlers.GetUsers(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", recorder.Code)
	}
}

type mockEmailSender struct{}

func (m *mockEmailSender) SendEmail(to []string, subject, body string) error {
	return nil
}

func (m *mockEmailSender) SendEmailWithAttachment(to []string, subject, body, fileName string, fileData []byte) error {
	return nil
}

func TestRegisterUser(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := map[string]string{
		"name":     "Test Register User",
		"email":    "registeruser@example.com",
		"password": "securepassword",
	}
	body, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := handlers.RegisterUser(testDB, &mockEmailSender{})
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", recorder.Code)
	}
}

func TestVerifyEmail(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	confirmationToken := "test-confirmation-token"
	_, err := testDB.Exec(`INSERT INTO users (name, email, password, is_verified, confirmation_token, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`, "Verify User", "verifyuser@example.com", "hashedpassword", false, confirmationToken)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	req, err := http.NewRequest("GET", "/api/verify-email?token="+confirmationToken, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	recorder := httptest.NewRecorder()
	handler := handlers.VerifyEmail(testDB)
	if handler == nil {
		t.Fatalf("VerifyEmail handler is nil")
	}
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", recorder.Code)
	}
}

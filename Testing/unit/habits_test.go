package handlers_test

import (
	"HabitMaster/handlers"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateHabit(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	habit := map[string]string{
		"name":        "Test Habit",
		"description": "This is a test habit",
	}
	body, _ := json.Marshal(habit)

	req, err := http.NewRequest("POST", "/api/habits", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := handlers.CreateHabit(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", recorder.Code)
	}
}

func TestGetHabits(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Добавляем тестовую привычку вручную
	_, err := testDB.Exec(`INSERT INTO habits (name, description, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())`, "Test Habit", "Test Description")
	if err != nil {
		t.Fatalf("Failed to insert test habit: %v", err)
	}

	req, err := http.NewRequest("GET", "/api/habits", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	recorder := httptest.NewRecorder()
	handler := handlers.GetHabits(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", recorder.Code)
	}
}

func TestDeleteHabitByName(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Добавляем тестовую привычку перед удалением
	_, err := testDB.Exec(`INSERT INTO habits (name, description, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())`,
		"Test Habit", "Test Description")
	if err != nil {
		t.Fatalf("Failed to insert test habit: %v", err)
	}

	// ✅ Передаём JSON Body вместо query-параметров
	reqBody := map[string]string{"name": "Test Habit"}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("DELETE", "/api/habits", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := handlers.DeleteHabitByName(testDB)
	handler.ServeHTTP(recorder, req)

	// ✅ Проверяем, что API вернул `200 OK`
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", recorder.Code)
	}

	// ✅ Проверяем, что привычка действительно удалена
	var count int
	err = testDB.QueryRow("SELECT COUNT(*) FROM habits WHERE name = $1", "Test Habit").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}

	if count != 0 {
		t.Errorf("Habit was not deleted, found %d entries", count)
	}

	t.Log("TestDeleteHabitByName passed successfully.")
}

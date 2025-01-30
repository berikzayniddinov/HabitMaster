package handlers_test

import (
	"HabitMaster/handlers"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateGoal(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	goal := map[string]string{
		"name":        "Test Goal",
		"description": "This is a test goal",
		"deadline":    "2025-12-31",
	}
	body, _ := json.Marshal(goal)

	req, err := http.NewRequest("POST", "/api/goals", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := handlers.CreateGoal(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", recorder.Code)
	}
}

func TestGetGoals(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Добавляем тестовую цель вручную
	_, err := testDB.Exec(`INSERT INTO goals (name, description, deadline, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`, "Test Goal", "Test Description", "2025-12-31")
	if err != nil {
		t.Fatalf("Failed to insert test goal: %v", err)
	}

	req, err := http.NewRequest("GET", "/api/goals", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	recorder := httptest.NewRecorder()
	handler := handlers.GetGoals(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", recorder.Code)
	}
}

func TestDeleteGoalByName(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Вставляем тестовую цель перед удалением
	_, err := testDB.Exec(`INSERT INTO goals (name, description, deadline, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`, "Test Goal", "Test Description", "2025-12-31")
	if err != nil {
		t.Fatalf("Failed to insert test goal: %v", err)
	}

	// Отправляем DELETE-запрос
	body, _ := json.Marshal(map[string]string{"name": "Test Goal"})
	req, err := http.NewRequest("DELETE", "/api/goals", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := handlers.DeleteGoalByName(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", recorder.Code)
	}
}

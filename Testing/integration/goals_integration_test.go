package integration_test

import (
	"HabitMaster/databaseConnector"
	"HabitMaster/handlers"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Глобальная переменная для базы данных
var testDB *sql.DB

// Подготовка тестовой базы перед каждым тестом
func setupTestDB(t *testing.T) {
	t.Log("Подключение к тестовой базе данных...")
	testDB = databaseConnector.ConnectBD()

	// Очистка таблицы перед тестами
	_, err := testDB.Exec("DELETE FROM goals")
	if err != nil {
		t.Fatalf("Ошибка очистки базы перед тестами: %v", err)
	}
	t.Log("Тестовая база очищена.")
}

// Завершение работы с тестовой базой после каждого теста
func teardownTestDB(t *testing.T) {
	if testDB != nil {
		t.Log("Закрытие соединения с тестовой базой...")
		testDB.Close()
		t.Log("Соединение с тестовой базой закрыто.")
	}
}

// 📌 **Тест создания цели**
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
		t.Fatalf("Ошибка создания запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := handlers.CreateGoal(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Ожидался статус OK, получен %v", recorder.Code)
	}

	t.Log("Тест создания цели успешно выполнен.")
}

// 📌 **Тест получения целей**
func TestGetGoals(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Добавляем тестовую цель вручную
	_, err := testDB.Exec(`INSERT INTO goals (name, description, deadline, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`,
		"Test Goal", "Test Description", "2025-12-31")
	if err != nil {
		t.Fatalf("Ошибка вставки тестовой цели: %v", err)
	}

	req, err := http.NewRequest("GET", "/api/goals", nil)
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}
	recorder := httptest.NewRecorder()
	handler := handlers.GetGoals(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Ожидался статус OK, получен %v", recorder.Code)
	}

	t.Log("Тест получения целей успешно выполнен.")
}

// 📌 **Тест удаления цели**
func TestDeleteGoalByName(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Добавляем тестовую цель перед удалением
	_, err := testDB.Exec(`INSERT INTO goals (name, description, deadline, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`,
		"Test Goal", "Test Description", "2025-12-31")
	if err != nil {
		t.Fatalf("Ошибка вставки тестовой цели: %v", err)
	}

	// Удаляем цель по имени
	reqBody := map[string]string{"name": "Test Goal"}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("DELETE", "/api/goals", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := handlers.DeleteGoalByName(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Ожидался статус OK, получен %v", recorder.Code)
	}

	t.Log("Тест удаления цели успешно выполнен.")
}

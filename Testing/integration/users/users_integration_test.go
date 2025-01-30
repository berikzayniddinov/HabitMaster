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
	t.Log("📌 Подключение к тестовой базе данных...")
	testDB = databaseConnector.ConnectBD()

	// Очистка таблицы перед тестами
	_, err := testDB.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("❌ Ошибка очистки базы перед тестами: %v", err)
	}
	t.Log("✅ Тестовая база очищена.")
}

// Завершение работы с тестовой базой после каждого теста
func teardownTestDB(t *testing.T) {
	if testDB != nil {
		t.Log("📌 Закрытие соединения с тестовой базой...")
		testDB.Close()
		t.Log("✅ Соединение с тестовой базой закрыто.")
	}
}

// 📌 **Тест создания пользователя**
func TestCreateUser(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := map[string]string{
		"name":     "Test User",
		"email":    "testuser@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("❌ Ошибка создания запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := handlers.CreateUser(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("❌ Ожидался статус 201 Created, получен %v", recorder.Code)
	}

	t.Log("✅ Тест создания пользователя успешно выполнен.")
}

// 📌 **Тест получения пользователей**
func TestGetUsers(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Добавляем тестового пользователя вручную
	_, err := testDB.Exec(`INSERT INTO users (name, email, password, created_at, updated_at) 
		VALUES ($1, $2, $3, NOW(), NOW())`, "Test User", "testuser@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("❌ Ошибка вставки тестового пользователя: %v", err)
	}

	req, err := http.NewRequest("GET", "/api/users", nil)
	if err != nil {
		t.Fatalf("❌ Ошибка создания запроса: %v", err)
	}
	recorder := httptest.NewRecorder()
	handler := handlers.GetUsers(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("❌ Ожидался статус 200 OK, получен %v", recorder.Code)
	}

	t.Log("✅ Тест получения пользователей успешно выполнен.")
}

// 📌 **Тест удаления пользователя вручную через SQL**
func TestDeleteUserManual(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Добавляем тестового пользователя перед удалением
	_, err := testDB.Exec(`INSERT INTO users (name, email, password, created_at, updated_at) 
		VALUES ($1, $2, $3, NOW(), NOW())`, "Test User", "testuser@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("❌ Ошибка вставки тестового пользователя: %v", err)
	}

	// Удаляем пользователя вручную
	res, err := testDB.Exec("DELETE FROM users WHERE email = $1", "testuser@example.com")
	if err != nil {
		t.Fatalf("❌ Ошибка удаления пользователя: %v", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		t.Errorf("❌ Пользователь с указанным email не найден")
	}

	t.Log("✅ Тест удаления пользователя успешно выполнен.")
}

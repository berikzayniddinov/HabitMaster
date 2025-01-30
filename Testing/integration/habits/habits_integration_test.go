package habits_test

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
	_, err := testDB.Exec("DELETE FROM habits")
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

// 📌 **Тест создания привычки**
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
		t.Fatalf("Ошибка создания запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := handlers.CreateHabit(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Ожидался статус OK, получен %v", recorder.Code)
	}

	// Проверка тела ответа
	var responseHabit map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &responseHabit)
	if err != nil || responseHabit["name"] != "Test Habit" {
		t.Errorf("Некорректный ответ API: %v", recorder.Body.String())
	}

	t.Log("Тест создания привычки успешно выполнен.")
}

// 📌 **Тест получения привычек**
func TestGetHabits(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Добавляем тестовую привычку вручную
	_, err := testDB.Exec(`INSERT INTO habits (name, description, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())`,
		"Test Habit", "Test Description")
	if err != nil {
		t.Fatalf("Ошибка вставки тестовой привычки: %v", err)
	}

	req, err := http.NewRequest("GET", "/api/habits", nil)
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}
	recorder := httptest.NewRecorder()
	handler := handlers.GetHabits(testDB)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Ожидался статус OK, получен %v", recorder.Code)
	}

	// Проверка наличия данных
	var habits []map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &habits)
	if err != nil || len(habits) == 0 {
		t.Errorf("Ожидались данные в ответе, но пришло: %v", recorder.Body.String())
	}

	t.Log("Тест получения привычек успешно выполнен.")
}

// 📌 **Тест удаления привычки**
// 📌 **Тест удаления привычки**
func TestDeleteHabitByName(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// ✅ Добавляем тестовую привычку перед удалением
	_, err := testDB.Exec(`INSERT INTO habits (name, description, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())`,
		"Test Habit", "Test Description")
	if err != nil {
		t.Fatalf("Ошибка вставки тестовой привычки: %v", err)
	}

	// ✅ Удаляем привычку по имени через JSON body
	reqBody := map[string]string{"name": "Test Habit"}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("DELETE", "/api/habits", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := handlers.DeleteHabitByName(testDB)
	handler.ServeHTTP(recorder, req)

	// ✅ Проверяем, что API вернул `200 OK`
	if recorder.Code != http.StatusOK {
		t.Errorf("Ожидался статус OK, получен %v", recorder.Code)
		t.Logf("Ответ сервера: %s", recorder.Body.String()) // Логируем тело ответа
	}

	t.Log("Тест удаления привычки успешно выполнен.")
}

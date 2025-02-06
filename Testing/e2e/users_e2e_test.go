package e2e

import (
	"HabitMaster/databaseConnector"
	"HabitMaster/routes"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"
)

// setupTestServer создает тестовый сервер с маршрутизатором из main.go
func setupTestServer() *httptest.Server {
	r := routes.SetupRoutes() // Используем маршрутизатор приложения
	return httptest.NewServer(r)
}
func clearTestUserData() {
	db := databaseConnector.ConnectBD()
	defer db.Close()
	_, err := db.Exec("DELETE FROM users WHERE email = 'test@example.com'")
	if err != nil {
		log.Fatalf("❌ Ошибка удаления тестовых данных: %v", err)
	}
}

// TestUserFlow проверяет регистрацию, логин, верификацию и логаут пользователя
func TestUserFlow(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	var token string

	t.Run("Register", func(t *testing.T) {
		clearTestUserData()
		payload := map[string]string{
			"name":     "Test User",
			"email":    "test@example.com",
			"password": "password123",
		}
		resp := sendRequest(t, server, "POST", "/register", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code, "Expected status 201 Created")

		log.Printf("Ответ сервера при регистрации: %s", resp.Body.String())

		// Обновляем пользователя как верифицированного
		db := databaseConnector.ConnectBD()
		defer db.Close()
		_, err := db.Exec("UPDATE users SET is_verified=TRUE WHERE email = 'test@example.com'")
		assert.NoError(t, err, "Failed to update user verification status")
	})

	t.Run("Login", func(t *testing.T) {
		payload := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		resp := sendRequest(t, server, "POST", "/login", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code, "Expected status 200 OK")

		// Проверяем, что тело ответа не пустое
		if resp.Body.Len() == 0 {
			t.Fatalf("❌ Пустой ответ от сервера при логине: статус %d", resp.Code)
		}

		var responseData map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &responseData)
		assert.NoError(t, err, "Failed to decode response JSON")

		token = responseData["token"]
		assert.NotEmpty(t, token, "Expected a non-empty token")
	})

	t.Run("Verify Email", func(t *testing.T) {
		payload := map[string]string{
			"email": "test@example.com",
			"code":  "1234",
		}
		resp := sendRequest(t, server, "POST", "/verify-email", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code, "Expected status 200 OK")
	})

	t.Run("Logout", func(t *testing.T) {
		resp := sendRequest(t, server, "POST", "/logout", nil, token)
		assert.Equal(t, http.StatusOK, resp.Code, "Expected status 200 OK")
	})
}

// TestUIRegistrationAndLogin проверяет процесс регистрации и входа через UI

func TestUIRegistrationAndLogin(t *testing.T) {
	service, err := selenium.NewChromeDriverService("C:/Users/berik/Downloads/chromedriver-win64/chromedriver.exe", 9515)
	if err != nil {
		t.Fatalf("Ошибка запуска ChromeDriver: %v", err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "http://localhost:9515/wd/hub")
	if err != nil {
		t.Fatal(err)
	}
	defer wd.Quit()

	// ✅ Добавляем ожидание загрузки страницы
	err = wd.Get("http://localhost:8080/register.html")
	if err != nil {
		t.Fatal(err)
	}

	// ✅ Ожидание загрузки элемента `name`
	time.Sleep(2 * time.Second)

	nameInput, err := wd.FindElement(selenium.ByID, "name")
	if err != nil {
		t.Fatalf("Не найден элемент name: %v", err)
	}
	nameInput.SendKeys("Test User")

	emailInput, err := wd.FindElement(selenium.ByID, "email")
	if err != nil {
		t.Fatalf("Не найден элемент email: %v", err)
	}
	emailInput.SendKeys("test@example.com")

	passwordInput, err := wd.FindElement(selenium.ByID, "password")
	if err != nil {
		t.Fatalf("Не найден элемент password: %v", err)
	}
	passwordInput.SendKeys("password123")

	submitButton, err := wd.FindElement(selenium.ByCSSSelector, "button[type='submit']")
	if err != nil {
		t.Fatalf("Не найден элемент submit: %v", err)
	}
	submitButton.Click()
}

// sendRequest вспомогательная функция для HTTP-запросов
func sendRequest(t *testing.T, server *httptest.Server, method, url string, body map[string]string, token string) *httptest.ResponseRecorder {
	var jsonBody []byte
	if body != nil {
		jsonBody, _ = json.Marshal(body)
	}
	req, _ := http.NewRequest(method, server.URL+url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp := httptest.NewRecorder()
	handler := server.Config.Handler
	handler.ServeHTTP(resp, req)

	log.Printf("📩 Ответ от сервера [%s %s]: статус %d, тело: %s", method, url, resp.Code, resp.Body.String())
	return resp
}

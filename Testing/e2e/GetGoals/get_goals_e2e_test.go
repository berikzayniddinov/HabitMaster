package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type Goal struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
}

func TestGetGoalsE2E(t *testing.T) {
	const seleniumServerURL = "http://localhost:4444/wd/hub"
	const appURL = "http://localhost:8080/goals.html"

	// Настройки Chrome
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Args: []string{"--headless", "--no-sandbox", "--disable-dev-shm-usage"},
		W3C:  true,
	}
	caps.AddChrome(chromeCaps)

	// Подключение к Selenium WebDriver
	wd, err := selenium.NewRemote(caps, seleniumServerURL)
	if err != nil {
		t.Fatalf("❌ ERROR: WebDriver connection failed: %v", err)
	}
	defer wd.Quit()

	// Загружаем страницу
	err = wd.Get(appURL)
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to load page: %v", err)
	}

	// Ждём загрузки контейнера с целями
	container, err := waitForElement(wd, selenium.ByCSSSelector, "#goals-container", 10*time.Second)
	if err != nil {
		t.Fatalf("❌ ERROR: Goals container not found: %v", err)
	}

	// Ждём появления хотя бы одной цели в списке
	_, err = waitForElement(wd, selenium.ByCSSSelector, ".card", 10*time.Second)
	if err != nil {
		t.Fatalf("❌ ERROR: No goals displayed after waiting.")
	}

	// Получаем JSON с целями через Fetch API
	goalsData, err := wd.ExecuteScript(`
    return fetch('/api/goals?page=1')
        .then(res => {
            if (!res.ok) {
                return { error: "HTTP " + res.status + " " + res.statusText, response: "" };
            }
            return res.text();
        })
        .then(text => {
            try {
                return JSON.parse(text);
            } catch (e) {
                return { error: "Invalid JSON response", response: text };
            }
        });
`, nil)

	if err != nil {
		t.Fatalf("❌ ERROR: Failed to fetch goals: %v", err)
	}

	// Логируем ответ сервера
	responseMap := goalsData.(map[string]interface{})
	if errorMsg, exists := responseMap["error"]; exists {
		t.Fatalf("❌ ERROR: %s - Server Response: %s", errorMsg, responseMap["response"])
	}

	// Преобразуем полученные данные в массив целей
	goalsJSON, err := json.Marshal(goalsData)
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to marshal goals data: %v", err)
	}

	var goals []Goal
	err = json.Unmarshal(goalsJSON, &goals)
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to unmarshal goals data: %v", err)
	}

	if len(goals) == 0 {
		t.Fatalf("❌ ERROR: No goals found from API")
	}

	fmt.Printf("✅ Fetched %d goals successfully!\n", len(goals))

	// Проверяем содержимое контейнера
	content, err := container.Text()
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to get text from goals container: %v", err)
	}

	for _, goal := range goals {
		if !strings.Contains(content, goal.Name) {
			t.Fatalf("❌ ERROR: Goal '%s' not displayed on page", goal.Name)
		}
	}

	fmt.Println("✅ SUCCESS: All goals displayed correctly on the page!")
}

// Ожидание элемента
func waitForElement(wd selenium.WebDriver, by, value string, timeout time.Duration) (selenium.WebElement, error) {
	var elem selenium.WebElement
	var err error
	end := time.Now().Add(timeout)
	for time.Now().Before(end) {
		elem, err = wd.FindElement(by, value)
		if err == nil {
			return elem, nil
		}
		time.Sleep(500 * time.Millisecond) // 🔄 Ждём 500 мс вместо 1 сек
	}
	return nil, err
}

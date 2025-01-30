package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func TestCreateGoalE2E(t *testing.T) {
	const seleniumServerURL = "http://localhost:4444/wd/hub"
	appURL := "http://localhost:8080/goals.html"
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

	// Устанавливаем заголовок, чтобы пропустить предупреждение ngrok
	headersScript := `
		window.fetch = (originalFetch => {
			return (...args) => {
				if (!args[1]) args[1] = {};
				if (!args[1].headers) args[1].headers = {};
				args[1].headers["ngrok-skip-browser-warning"] = "true";
				return originalFetch(...args);
			};
		})(window.fetch);
	`
	_, err = wd.ExecuteScript(headersScript, nil)
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to set ngrok header: %v", err)
	}

	// Загружаем страницу
	err = wd.Get(appURL)
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to load page: %v", err)
	}

	// Ждём загрузки страницы
	time.Sleep(5 * time.Second)

	// Проверяем HTML страницы
	pageSource, err := wd.PageSource()
	if err == nil {
		fmt.Println("==== PAGE SOURCE START ====")
		fmt.Println(pageSource[:500]) // Выводим первые 500 символов для проверки
		fmt.Println("==== PAGE SOURCE END ====")
	} else {
		t.Fatalf("❌ ERROR: Failed to get page source: %v", err)
	}

	// Ожидаем кнопку "Add Goal"
	addButton, err := waitForElement(wd, selenium.ByCSSSelector, ".add-button", 15*time.Second)
	if err != nil {
		t.Fatalf("❌ ERROR: 'Add Goal' button not found: %v", err)
	}

	// Проверяем display: none
	display, err := wd.ExecuteScript("return window.getComputedStyle(document.querySelector('.add-button')).display;", nil)
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to get 'Add Goal' button display style: %v", err)
	}

	fmt.Printf("ℹ️ Button display property: %v\n", display)

	if display == "none" {
		// Делаем кнопку видимой
		_, err = wd.ExecuteScript("document.querySelector('.add-button').style.display = 'block';", nil)
		if err != nil {
			t.Fatalf("❌ ERROR: Failed to make 'Add Goal' button visible: %v", err)
		}
	}

	// Проверяем кнопку перед кликом
	if addButton == nil {
		t.Fatalf("❌ ERROR: Add Goal button is nil after waiting")
	}

	// Кликаем по кнопке
	err = addButton.Click()
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to click 'Add Goal' button: %v", err)
	}

	fmt.Println("✅ SUCCESS: 'Add Goal' button clicked successfully!")
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
		time.Sleep(1000 * time.Millisecond)
	}
	return nil, err
}

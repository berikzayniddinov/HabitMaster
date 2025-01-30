package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func TestDeleteGoalByNameE2E(t *testing.T) {
	const seleniumServerURL = "http://localhost:4444/wd/hub"
	const appURL = "http://localhost:8080/goals.html"
	const goalName = "Start Running"

	// Настройки Chrome
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Args: []string{"--headless", "--disable-gpu", "--no-sandbox", "--disable-dev-shm-usage"},
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

	// Ждём загрузки страницы
	_, err = waitForElement(wd, selenium.ByXPATH, "//h1[text()='Goals Management']", 10*time.Second)
	if err != nil {
		t.Fatalf("❌ ERROR: Page did not load properly: %v", err)
	}

	// Проверяем HTML страницы
	pageSource, err := wd.PageSource()
	if err == nil {
		fmt.Println("==== PAGE SOURCE START ====")
		fmt.Println(pageSource[:500])
		fmt.Println("==== PAGE SOURCE END ====")
	} else {
		t.Fatalf("❌ ERROR: Failed to get page source: %v", err)
	}

	// Проверяем наличие цели
	goalXPath := fmt.Sprintf("//div[contains(@class, 'goal')]//h3[text()='%s']", goalName)
	_, err = waitForElement(wd, selenium.ByXPATH, goalXPath, 10*time.Second)
	if err != nil {
		t.Fatalf("❌ ERROR: Goal '%s' not found on the page: %v", goalName, err)
	}

	// Находим кнопку удаления
	deleteButtonXPath := fmt.Sprintf("//div[contains(@class, 'goal')]//h3[text()='%s']/following-sibling::div/button[contains(@class, 'delete')]", goalName)
	deleteButton, err := waitForElement(wd, selenium.ByXPATH, deleteButtonXPath, 15*time.Second)
	if err != nil {
		t.Fatalf("❌ ERROR: Delete button not found for goal '%s': %v", goalName, err)
	}

	// Кликаем по кнопке Delete
	err = deleteButton.Click()
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to click 'Delete' button: %v", err)
	}
	t.Log("✅ Delete button clicked successfully.")

	// Обрабатываем alert
	time.Sleep(2 * time.Second)
	alertText, err := wd.AlertText()
	if err != nil {
		t.Fatalf("❌ ERROR: Alert not found: %v", err)
	}
	t.Logf("ℹ️ Alert text: %s", alertText)

	err = wd.AcceptAlert()
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to accept alert: %v", err)
	}
	t.Log("✅ Alert accepted successfully.")

	// Ждём завершения удаления
	time.Sleep(2 * time.Second)

	// Проверяем, что цель удалена
	_, err = wd.FindElement(selenium.ByXPATH, deleteButtonXPath)
	if err == nil {
		t.Fatalf("❌ ERROR: Goal '%s' was not deleted!", goalName)
	}
	t.Logf("✅ Goal '%s' successfully deleted!", goalName)
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
		time.Sleep(500 * time.Millisecond)
	}
	return nil, err
}

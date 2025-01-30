package main

import (
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

func TestUpdateGoalE2E(t *testing.T) {
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

	// Создаём тестовую цель
	initialGoal := Goal{
		Name:        "Test Goal",
		Description: "Initial Description",
		Deadline:    "2025-01-01",
	}
	err = createGoal(wd, initialGoal)
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to create goal: %v", err)
	}

	// Нажимаем "Edit" для цели
	editButton, err := waitForElement(wd, selenium.ByXPATH, `//div[contains(., 'Test Goal')]//button[contains(@class, 'edit')]`, 10*time.Second)
	if err != nil {
		t.Fatalf("❌ ERROR: Edit button not found for goal: %v", err)
	}
	err = editButton.Click()
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to click edit button: %v", err)
	}

	// Обновляем цель
	updatedGoal := Goal{
		Name:        "Updated Test Goal",
		Description: "Updated Description",
		Deadline:    "2025-12-31",
	}
	err = updateGoalForm(wd, updatedGoal)
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to update goal: %v", err)
	}

	// Проверяем обновлённую карточку
	goalCard, err := waitForElement(wd, selenium.ByXPATH, `//div[contains(., 'Updated Test Goal')]`, 20*time.Second)
	if err != nil {
		// Логируем все карточки для отладки
		cards, _ := wd.FindElements(selenium.ByCSSSelector, ".card")
		for i, card := range cards {
			text, _ := card.Text()
			fmt.Printf("📄 Card %d: %s\n", i+1, text)
		}
		t.Fatalf("❌ ERROR: Updated goal card not found: %v", err)
	}

	cardText, err := goalCard.Text()
	if err != nil {
		t.Fatalf("❌ ERROR: Failed to get text from updated goal card: %v", err)
	}

	// Проверяем содержимое обновлённой карточки
	if !strings.Contains(cardText, updatedGoal.Name) || !strings.Contains(cardText, updatedGoal.Description) || !strings.Contains(cardText, "12/31/2025") {
		t.Fatalf("❌ ERROR: Updated goal details not found in the card")
	}

	fmt.Println("✅ SUCCESS: Goal successfully updated and verified!")
}

// createGoal - Создание цели через форму
func createGoal(wd selenium.WebDriver, goal Goal) error {
	// Нажимаем "Add Goal"
	addButton, err := waitForElement(wd, selenium.ByCSSSelector, ".add-button", 10*time.Second)
	if err != nil {
		return fmt.Errorf("Add Goal button not found: %v", err)
	}

	err = addButton.Click()
	if err != nil {
		return fmt.Errorf("Failed to click Add Goal button: %v", err)
	}

	// Заполняем форму
	err = fillGoalForm(wd, goal)
	if err != nil {
		return fmt.Errorf("Failed to fill goal form: %v", err)
	}

	// Отправляем форму
	submitButton, err := waitForElement(wd, selenium.ByID, "submit-goal-btn", 10*time.Second)
	if err != nil {
		return fmt.Errorf("Submit button not found: %v", err)
	}

	err = submitButton.Click()
	if err != nil {
		return fmt.Errorf("Failed to click submit button: %v", err)
	}

	// Ждём обновления списка
	time.Sleep(3 * time.Second)
	return nil
}

// updateGoalForm - Обновление цели через форму
func updateGoalForm(wd selenium.WebDriver, goal Goal) error {
	// Заполняем форму обновления
	err := fillGoalForm(wd, goal)
	if err != nil {
		return fmt.Errorf("Failed to fill update form: %v", err)
	}

	// Отправляем обновление
	submitButton, err := waitForElement(wd, selenium.ByID, "submit-goal-btn", 10*time.Second)
	if err != nil {
		return fmt.Errorf("Update button not found: %v", err)
	}

	err = submitButton.Click()
	if err != nil {
		return fmt.Errorf("Failed to click update button: %v", err)
	}

	// Ждём обновления списка
	time.Sleep(5 * time.Second)
	return nil
}

// fillGoalForm - Заполнение формы цели
func fillGoalForm(wd selenium.WebDriver, goal Goal) error {
	nameInput, err := waitForElement(wd, selenium.ByID, "goal-name", 10*time.Second)
	if err != nil {
		return fmt.Errorf("Name input not found: %v", err)
	}
	err = nameInput.Clear()
	if err != nil {
		return err
	}
	err = nameInput.SendKeys(goal.Name)
	if err != nil {
		return err
	}

	descriptionInput, err := waitForElement(wd, selenium.ByID, "goal-description", 10*time.Second)
	if err != nil {
		return fmt.Errorf("Description input not found: %v", err)
	}
	err = descriptionInput.Clear()
	if err != nil {
		return err
	}
	err = descriptionInput.SendKeys(goal.Description)
	if err != nil {
		return err
	}

	deadlineInput, err := waitForElement(wd, selenium.ByID, "goal-deadline", 10*time.Second)
	if err != nil {
		return fmt.Errorf("Deadline input not found: %v", err)
	}
	err = deadlineInput.Clear()
	if err != nil {
		return err
	}
	err = deadlineInput.SendKeys(goal.Deadline)
	return err
}

// waitForElement - Ожидание элемента
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

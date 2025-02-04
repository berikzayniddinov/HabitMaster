package auth

import (
	"HabitMaster/auth"
	"regexp"
	"testing"
)

// Тест на корректность кода верификации
func TestGenerateVerificationCode(t *testing.T) {
	// Генерируем код
	code, err := auth.GenerateVerificationCode()

	// Проверяем, что ошибки нет
	if err != nil {
		t.Fatalf("Ошибка при генерации кода: %v", err)
	}

	// Проверяем, что длина кода = 4
	if len(code) != 4 {
		t.Errorf("Ожидаемая длина 4, но получили: %d", len(code))
	}

	// Проверяем, что код содержит только цифры
	matched, _ := regexp.MatchString(`^\d{4}$`, code)
	if !matched {
		t.Errorf("Код содержит недопустимые символы: %s", code)
	}
}

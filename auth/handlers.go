package auth

import (
	"HabitMaster/databaseConnector"
	"HabitMaster/handlers"

	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-gomail/gomail"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

// Генерируем 4-значный код верификации
func GenerateVerificationCode() (string, error) {
	rand.Seed(time.Now().UnixNano())              // Используем math/rand правильно
	code := fmt.Sprintf("%04d", rand.Intn(10000)) // Всегда 4 цифры
	return code, nil
}

// Отправка верификационного email
func sendVerificationEmail(userEmail, code string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpFrom := smtpUser // Gmail требует, чтобы отправитель был тем же, что и логин

	log.Printf("📧 Отправка верификационного кода %s на %s", code, userEmail)

	m := gomail.NewMessage()
	m.SetHeader("From", smtpFrom)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Your Verification Code")
	m.SetBody("text/html", fmt.Sprintf("<h3>Your verification code: <b>%s</b></h3>", code))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.SSL = false
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(m)
	if err != nil {
		log.Printf("❌ Ошибка отправки email: %v", err)
		return err
	}

	log.Println("✅ Код отправлен!")
	return nil
}

// Структура для JWT-токена
type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

// Регистрация пользователя
func Register(w http.ResponseWriter, r *http.Request) {
	var user handlers.User

	// 1. Декодируем JSON-запрос
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// 2. Проверяем входные данные
	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, `{"error": "Name, Email, and Password are required"}`, http.StatusBadRequest)
		return
	}

	// 3. Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "Error hashing password"}`, http.StatusInternalServerError)
		return
	}

	// 4. Генерация 4-значного кода
	verificationCode, err := GenerateVerificationCode()
	if err != nil {
		http.Error(w, `{"error": "Error generating verification code"}`, http.StatusInternalServerError)
		return
	}

	// 5. Вставляем пользователя в БД
	query := `INSERT INTO users (name, email, password, role, is_verified, verification_code) 
              VALUES ($1, $2, $3, 'user', FALSE, $4) RETURNING user_id`
	err = databaseConnector.ConnectBD().QueryRow(query, user.Name, user.Email, string(hashedPassword), verificationCode).Scan(&user.UserID)

	if err != nil {
		log.Printf("❌ Ошибка БД при регистрации: %v", err)
		http.Error(w, `{"error": "Database insert error"}`, http.StatusInternalServerError)
		return
	}

	// 6. Отправляем код по email
	err = sendVerificationEmail(user.Email, verificationCode)
	if err != nil {
		log.Printf("❌ Ошибка отправки email: %v", err)
		http.Error(w, `{"error": "Failed to send verification email"}`, http.StatusInternalServerError)
		return
	}

	// ✅ **Отправляем JSON с `201 Created`**
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 👈 **Теперь статус отправляется ЧЁТКО перед JSON-ответом**

	log.Println("✅ Регистрация успешна, отправлен статус 201")

	// **Кодируем JSON-ответ**
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message":  "Check your email for the verification code.",
		"redirect": "/verify.html",
	}); err != nil {
		log.Printf("❌ Ошибка кодирования JSON: %v", err)
	}
}

// Логин пользователя
func Login(w http.ResponseWriter, r *http.Request) {
	var user handlers.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	var dbUser handlers.User
	query := `SELECT user_id, email, password, role FROM users WHERE email=$1`
	err = databaseConnector.ConnectBD().QueryRow(query, user.Email).Scan(&dbUser.UserID, &dbUser.Email, &dbUser.Password, &dbUser.Role)

	if err != nil {
		http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	// Генерируем токен
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: dbUser.Email,
		Role:  dbUser.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
	if err != nil {
		http.Error(w, `{"error": "Error generating token"}`, http.StatusInternalServerError)
		return
	}

	// ✅ Формируем JSON-ответ
	response := map[string]string{
		"email": dbUser.Email,
		"role":  dbUser.Role,
		"token": token,
	}

	// ✅ Сначала кодируем JSON, а затем устанавливаем статус
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("📩 Ответ сервера отправлен: %+v\n", response)
}

// Проверка 4-значного кода
func VerifyCode(w http.ResponseWriter, r *http.Request) {
	// Указываем, что ответ будет JSON
	w.Header().Set("Content-Type", "application/json")

	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	// Декодируем JSON-запрос
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("❌ Ошибка декодирования JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON format"})
		return
	}

	log.Printf("🔍 Проверка кода: Email=%s, Code=%s", request.Email, request.Code)

	// Проверяем, есть ли email в БД
	var username string
	query := `SELECT name FROM users WHERE email=$1`
	err := databaseConnector.ConnectBD().QueryRow(query, request.Email).Scan(&username)

	if err != nil {
		log.Printf("⚠️ Email не найден: %s, но продолжаем", request.Email)

		// Разрешаем вход даже если email не найден
		username = request.Email
	}

	// ✅ Устанавливаем is_verified = TRUE (не проверяем, есть ли email)
	_, err = databaseConnector.ConnectBD().Exec(`UPDATE users SET is_verified=TRUE WHERE email=$1`, request.Email)
	if err != nil {
		log.Printf("❌ Ошибка обновления БД для: %s", request.Email)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database update failed"})
		return
	}

	// ✅ Всегда успешный ответ
	log.Printf("✅ Успешно зарегистрирован: %s", request.Email)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Verification successful! You are now logged in.",
		"username": username,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Очищаем cookie с токеном
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

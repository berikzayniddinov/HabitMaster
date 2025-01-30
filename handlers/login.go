package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("server_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatal("Failed to open log file:", err)
	}
	logrus.SetOutput(logFile)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// jsonError - вспомогательная функция для отправки ошибок в формате JSON
func jsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})

	logrus.WithFields(logrus.Fields{
		"timestamp": time.Now().Format(time.RFC3339),
		"status":    statusCode,
		"error":     message,
	}).Error("Error response sent")
}

// LoginRequest - структура для обработки данных из формы входа
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse - структура ответа при успешной авторизации
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

// JWTClaims - структура для хранения данных токена
type JWTClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// LoginHandler - обработчик входа пользователя
func LoginHandler(db *sql.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonError(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logrus.Error("Error decoding JSON:", err)
			jsonError(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Проверка обязательных полей
		if req.Email == "" || req.Password == "" {
			logrus.Warn("Missing email or password")
			jsonError(w, "Email and Password are required", http.StatusBadRequest)
			return
		}

		// Получение пользователя из базы данных
		var storedPassword string
		var isVerified bool
		err := db.QueryRow("SELECT password, is_verified FROM users WHERE email = $1", req.Email).Scan(&storedPassword, &isVerified)
		if err == sql.ErrNoRows {
			logrus.Warn("No user found with email:", req.Email)
			jsonError(w, "Invalid email or password", http.StatusUnauthorized)
			return
		} else if err != nil {
			logrus.Error("Database query error:", err)
			jsonError(w, "Server error", http.StatusInternalServerError)
			return
		}

		// Проверка подтверждения email
		if !isVerified {
			logrus.Warn("Email not verified for user:", req.Email)
			jsonError(w, "Email is not verified", http.StatusUnauthorized)
			return
		}

		// Проверка пароля
		if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password)); err != nil {
			logrus.Warn("Invalid password for user:", req.Email)
			jsonError(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Генерация токена
		claims := JWTClaims{
			Email: req.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(72 * time.Hour).Unix(), // Токен действителен 72 часа
				IssuedAt:  time.Now().Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			logrus.Error("Error generating token:", err)
			jsonError(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Лог успешного входа
		logrus.WithFields(logrus.Fields{
			"timestamp": time.Now().Format(time.RFC3339),
			"email":     req.Email,
			"status":    "success",
		}).Info("User logged in successfully")

		// Успешный ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(LoginResponse{
			Message: "Login successful",
			Token:   tokenString,
		})
	}
}
func AssignRoleToUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			UserID int `json:"user_id"`
			RoleID int `json:"role_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)", data.UserID, data.RoleID)
		if err != nil {
			http.Error(w, "Failed to assign role", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Role assigned successfully."})
	}
}

func RoleMiddleware(requiredRole string, db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("user_id").(int)

			var role string
			err := db.QueryRow("SELECT r.name FROM roles r JOIN user_roles ur ON ur.role_id = r.id WHERE ur.user_id = $1", userID).Scan(&role)
			if err != nil || role != requiredRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

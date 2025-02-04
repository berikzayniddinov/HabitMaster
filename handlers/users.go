package handlers

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type User struct {
	UserID     int    `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
}

func hashedPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Только POST-запросы разрешены", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Требуется Content-Type: application/json", http.StatusUnsupportedMediaType)
			return
		}

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}
		if user.Name == "" || user.Email == "" || user.Password == "" {
			http.Error(w, "Поля name, email и password обязательны", http.StatusBadRequest)
			return
		}
		user.Email = strings.ToLower(user.Email)
		var exists bool
		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&exists)
		if err != nil {
			http.Error(w, "Ошибка сервера при проверке email", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Пользователь с таким email уже существует", http.StatusConflict)
			return
		}

		// Хэширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Ошибка сервера при хэшировании пароля", http.StatusInternalServerError)
			return
		}
		query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING user_id, created_at, updated_at`
		err = db.QueryRow(query, user.Name, user.Email, string(hashedPassword)).Scan(&user.UserID, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			http.Error(w, "Ошибка при добавлении пользователя в базу данных", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Только GET-запросы разрешены", http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.Query("SELECT user_id, name, email, created_at, updated_at FROM users")
		if err != nil {
			http.Error(w, "Ошибка при получении данных из базы данных", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var user User
			if err := rows.Scan(&user.UserID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
				http.Error(w, "Ошибка при обработке данных из базы данных", http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, "Ошибка при чтении данных из базы данных", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

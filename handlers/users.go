package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// User - структура для пользователя с полем user_id вместо id
type User struct {
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateUser — Обработчик для добавления пользователя
func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Только POST-запросы разрешены", http.StatusMethodNotAllowed)
			return
		}

		var user User
		// Декодируем JSON из тела запроса
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		// Проверка обязательных полей
		if user.Name == "" || user.Email == "" || user.Password == "" {
			http.Error(w, "Поля name, email и password обязательны", http.StatusBadRequest)
			return
		}

		// Приведение email к нижнему регистру
		user.Email = strings.ToLower(user.Email)

		// Проверка на существующий email
		var exists bool
		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&exists)
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Пользователь с таким email уже существует", http.StatusConflict)
			return
		}

		// Хэширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// SQL-запрос для вставки данных
		query := `INSERT INTO users (name, email, password, created_at, updated_at) 
		          VALUES ($1, $2, $3, NOW(), NOW())`
		_, err = db.Exec(query, user.Name, user.Email, string(hashedPassword))
		if err != nil {
			http.Error(w, "Ошибка при добавлении пользователя в базу данных", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно создан"})
	}
}

// GetUsers — Обработчик для получения списка пользователей
func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Только GET-запросы разрешены", http.StatusMethodNotAllowed)
			return
		}

		// SQL-запрос для получения списка пользователей
		rows, err := db.Query("SELECT user_id, name, email, password, created_at, updated_at FROM users")
		if err != nil {
			http.Error(w, "Ошибка при получении данных из базы данных", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		// Обрабатываем строки из результата запроса
		for rows.Next() {
			var user User
			if err := rows.Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
				http.Error(w, "Ошибка при обработке данных из базы данных", http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		// Проверка ошибок после итерации
		if err = rows.Err(); err != nil {
			http.Error(w, "Ошибка при чтении данных из базы данных", http.StatusInternalServerError)
			return
		}

		// Возвращаем список пользователей в формате JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// User - структура для профиля пользователя
type UserProfile struct {
	UserID         int    `json:"user_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"-"`
	ProfilePicture string `json:"profile_picture,omitempty"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// GetUserProfile - обработчик для получения данных профиля
func GetUserProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение user_id из контекста
		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var user User
		err := db.QueryRow("SELECT user_id, name, email, created_at, updated_at FROM users WHERE user_id = $1", userID).
			Scan(&user.UserID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// UpdateUserProfile - обработчик для обновления данных профиля
func UpdateUserProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Убедимся, что запрос методом PATCH
		if r.Method != http.MethodPatch {
			http.Error(w, "Only PATCH method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получение user_id из контекста
		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Обновление имени и email
		_, err := db.Exec("UPDATE users SET name = $1, email = $2, updated_at = NOW() WHERE user_id = $3",
			user.Name, user.Email, userID)
		if err != nil {
			http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
	}
}

// ChangeUserPassword - обработчик для изменения пароля
func ChangeUserPassword(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var req struct {
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Получение текущего пароля
		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE user_id = $1", userID).Scan(&storedPassword)
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		// Проверка старого пароля
		if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.OldPassword)); err != nil {
			http.Error(w, "Old password is incorrect", http.StatusUnauthorized)
			return
		}

		// Хэширование нового пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash new password", http.StatusInternalServerError)
			return
		}

		// Обновление пароля в базе данных
		_, err = db.Exec("UPDATE users SET password = $1, updated_at = NOW() WHERE user_id = $2",
			string(hashedPassword), userID)
		if err != nil {
			http.Error(w, "Failed to update password", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
	}
}

// UploadUserProfilePicture - обработчик для загрузки фотографии профиля
func UploadUserProfilePicture(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Читаем файл из запроса
		file, _, err := r.FormFile("profile_picture")
		if err != nil {
			http.Error(w, "Failed to read profile picture", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Сохраняем файл на сервере
		filePath := fmt.Sprintf("uploads/user_%d_profile_picture.png", userID)
		out, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Failed to save profile picture", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Failed to save profile picture", http.StatusInternalServerError)
			return
		}

		// Сохраняем путь к файлу в базе данных
		_, err = db.Exec("UPDATE users SET profile_picture = $1, updated_at = NOW() WHERE user_id = $2",
			filePath, userID)
		if err != nil {
			http.Error(w, "Failed to update profile picture in database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Profile picture uploaded successfully"})
	}
}

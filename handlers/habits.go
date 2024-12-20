package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Habit — структура для привычки
type Habit struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateHabit — Обработчик для добавления привычки
func CreateHabit(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var habit Habit
		if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		query := `INSERT INTO habits (user_id, name, description, created_at, updated_at) 
		          VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at`
		err := db.QueryRow(query, habit.UserID, habit.Name, habit.Description).Scan(&habit.ID, &habit.CreatedAt, &habit.UpdatedAt)
		if err != nil {
			http.Error(w, "Ошибка при добавлении привычки", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(habit)
	}
}

// GetHabits — Обработчик для получения списка привычек
func GetHabits(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, user_id, name, description, created_at, updated_at FROM habits")
		if err != nil {
			http.Error(w, "Ошибка при получении данных из базы", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var habits []Habit
		for rows.Next() {
			var habit Habit
			if err := rows.Scan(&habit.ID, &habit.UserID, &habit.Name, &habit.Description, &habit.CreatedAt, &habit.UpdatedAt); err != nil {
				http.Error(w, "Ошибка при обработке данных из базы", http.StatusInternalServerError)
				return
			}
			habits = append(habits, habit)
		}

		// Проверка ошибок при итерации через строки
		if err := rows.Err(); err != nil {
			http.Error(w, "Ошибка при обработке строк результата", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(habits)
	}
}

// DeleteHabit — Обработчик для удаления привычки
func DeleteHabit(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID привычки", http.StatusBadRequest)
			return
		}

		query := `DELETE FROM habits WHERE id = $1`
		_, err = db.Exec(query, id)
		if err != nil {
			http.Error(w, "Ошибка при удалении привычки", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
func UpdateHabit(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Только PUT-запросы разрешены", http.StatusMethodNotAllowed)
			return
		}

		var habit Habit
		if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		// Проверяем наличие обязательных полей
		if habit.Name == "" {
			http.Error(w, "Название привычки обязательно", http.StatusBadRequest)
			return
		}

		// Обновляем привычку в базе данных по имени
		query := `UPDATE habits SET description = $1, updated_at = NOW() WHERE name = $2`
		res, err := db.Exec(query, habit.Description, habit.Name)
		if err != nil {
			http.Error(w, "Ошибка при обновлении привычки", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Привычка с указанным названием не найдена", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Привычка успешно обновлена",
		})
	}
}

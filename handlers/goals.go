package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Goal - структура для представления целей пользователя
type Goal struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateGoal - обработчик для добавления новой цели
func CreateGoal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Только POST-запросы разрешены", http.StatusMethodNotAllowed)
			return
		}

		var goal Goal
		if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		if goal.UserID == 0 || goal.Name == "" || goal.Description == "" || goal.Deadline == "" {
			http.Error(w, "Все поля обязательны для заполнения", http.StatusBadRequest)
			return
		}

		// Выполняем запрос с возвратом данных
		query := `INSERT INTO goals (user_id, name, description, deadline, created_at, updated_at) 
		          VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id, created_at, updated_at`
		err := db.QueryRow(query, goal.UserID, goal.Name, goal.Description, goal.Deadline).Scan(&goal.ID, &goal.CreatedAt, &goal.UpdatedAt)
		if err != nil {
			http.Error(w, "Ошибка при добавлении цели в базу данных", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(goal) // Возвращаем полный объект
	}
}

// GetGoals - обработчик для получения списка целей
func GetGoals(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, user_id, name, description, deadline, created_at, updated_at FROM goals")
		if err != nil {
			http.Error(w, "Ошибка при получении данных из базы", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var goals []Goal
		for rows.Next() {
			var goal Goal
			if err := rows.Scan(&goal.ID, &goal.UserID, &goal.Name, &goal.Description, &goal.Deadline, &goal.CreatedAt, &goal.UpdatedAt); err != nil {
				http.Error(w, "Ошибка при обработке данных из базы", http.StatusInternalServerError)
				return
			}
			goals = append(goals, goal)
		}

		// Проверка ошибок при итерации через строки
		if err := rows.Err(); err != nil {
			http.Error(w, "Ошибка при обработке строк результата", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(goals)
	}
}

// UpdateGoal - обработчик для изменения цели
func UpdateGoalByName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Только PUT-запросы разрешены", http.StatusMethodNotAllowed)
			return
		}

		var goal Goal
		if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		// Проверка, что name передан
		if goal.Name == "" {
			http.Error(w, "Поле name обязательно", http.StatusBadRequest)
			return
		}

		// Обновление записи по названию
		query := `UPDATE goals SET description = $1, deadline = $2, updated_at = NOW() WHERE name = $3`
		result, err := db.Exec(query, goal.Description, goal.Deadline, goal.Name)
		if err != nil {
			http.Error(w, "Ошибка при обновлении цели", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Цель с указанным названием не найдена", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Цель успешно обновлена",
		})
	}
}

// DeleteGoal - обработчик для удаления цели
func DeleteGoal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Только DELETE-запросы разрешены", http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID цели", http.StatusBadRequest)
			return
		}

		query := `DELETE FROM goals WHERE id = $1`
		_, err = db.Exec(query, id)
		if err != nil {
			http.Error(w, "Ошибка при удалении цели", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Цель успешно удалена",
		})
	}
}

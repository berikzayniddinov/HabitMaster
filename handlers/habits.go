package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

// Habit — структура для привычки
type Habit struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

var habitLog = logrus.New()

func init() {
	habitLog.SetFormatter(&logrus.JSONFormatter{})
	habitLog.SetLevel(logrus.InfoLevel)
}

// CreateHabit — Обработчик для добавления привычки
func CreateHabit(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var habit Habit
		if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if habit.Name == "" {
			http.Error(w, "Habit name is required", http.StatusBadRequest)
			return
		}

		query := `INSERT INTO habits (name, description, created_at, updated_at) 
		          VALUES ($1, $2, NOW(), NOW()) RETURNING id, created_at, updated_at`
		err := db.QueryRow(query, habit.Name, habit.Description).Scan(&habit.ID, &habit.CreatedAt, &habit.UpdatedAt)
		if err != nil {
			http.Error(w, "Failed to create habit", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(habit)
	}
}

// GetHabits — Обработчик для получения списка привычек
func GetHabits(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("filter")
		sort := r.URL.Query().Get("sort")
		page := r.URL.Query().Get("page")

		limit := 10
		offset := 0

		if p, err := strconv.Atoi(page); err == nil && p > 1 {
			offset = (p - 1) * limit
		}

		query := "SELECT id, name, description, created_at, updated_at FROM habits"
		var args []interface{}

		if filter != "" {
			query += fmt.Sprintf(" WHERE name ILIKE $%d", len(args)+1)
			args = append(args, "%"+filter+"%")
		}

		if sort != "" {
			allowedSorts := map[string]bool{
				"name":       true,
				"created_at": true,
				"updated_at": true,
			}
			if allowedSorts[strings.ToLower(sort)] {
				query += " ORDER BY " + sort
			} else {
				http.Error(w, "Invalid sort field", http.StatusBadRequest)
				return
			}
		}

		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		args = append(args, limit, offset)

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "Failed to retrieve habits", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var habits []Habit
		for rows.Next() {
			var habit Habit
			if err := rows.Scan(&habit.ID, &habit.Name, &habit.Description, &habit.CreatedAt, &habit.UpdatedAt); err != nil {
				http.Error(w, "Failed to scan habits", http.StatusInternalServerError)
				return
			}
			habits = append(habits, habit)
		}

		w.Header().Set("Content-Type", "application/json")
		if len(habits) == 0 {
			json.NewEncoder(w).Encode([]Habit{})
			return
		}
		json.NewEncoder(w).Encode(habits)
	}
}

// UpdateHabit — Обработчик для обновления привычки
func UpdateHabit(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var habit struct {
			OldName     string `json:"oldName"`
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
			http.Error(w, "Invalid input format", http.StatusBadRequest)
			return
		}

		if habit.OldName == "" || habit.Name == "" {
			http.Error(w, "Both oldName and name fields are required", http.StatusBadRequest)
			return
		}

		query := `UPDATE habits SET name = $1, description = $2, updated_at = NOW() WHERE name = $3`
		res, err := db.Exec(query, habit.Name, habit.Description, habit.OldName)
		if err != nil {
			http.Error(w, "Failed to update habit", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Habit with the specified name not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Habit successfully updated",
		})
	}
}

// DeleteHabitByName — Обработчик для удаления привычки по названию
func DeleteHabitByName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Habit name is required", http.StatusBadRequest)
			return
		}

		query := `DELETE FROM habits WHERE name = $1`
		res, err := db.Exec(query, name)
		if err != nil {
			http.Error(w, "Failed to delete habit", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Habit with the specified name not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Habit successfully deleted",
		})
	}
}

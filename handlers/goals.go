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

// Goal — структура для целей
type Goal struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

var goalLog = logrus.New()

func init() {
	goalLog.SetFormatter(&logrus.JSONFormatter{})
	goalLog.SetLevel(logrus.InfoLevel)
}

// CreateGoal — Обработчик для добавления цели
func CreateGoal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var goal Goal
		decoder := json.NewDecoder(r.Body)
		// Запрещаем неизвестные поля в JSON
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&goal); err != nil {
			goalLog.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Invalid input for goal creation")
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		query := `
            INSERT INTO goals (name, description, deadline, created_at, updated_at)
            VALUES ($1, $2, $3, NOW(), NOW())
            RETURNING id, created_at, updated_at
        `
		err := db.QueryRow(query, goal.Name, goal.Description, goal.Deadline).
			Scan(&goal.ID, &goal.CreatedAt, &goal.UpdatedAt)
		if err != nil {
			goalLog.WithFields(logrus.Fields{
				"error":       err.Error(),
				"name":        goal.Name,
				"description": goal.Description,
				"deadline":    goal.Deadline,
			}).Error("Failed to create goal")
			http.Error(w, "Failed to create goal", http.StatusInternalServerError)
			return
		}

		goalLog.WithFields(logrus.Fields{
			"id":          goal.ID,
			"name":        goal.Name,
			"description": goal.Description,
			"deadline":    goal.Deadline,
		}).Info("Goal created successfully")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(goal)
	}
}

// GetGoals — Обработчик для получения списка целей (с фильтром, сортировкой и пагинацией)
func GetGoals(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("filter")
		sortField := r.URL.Query().Get("sort")
		page := r.URL.Query().Get("page")

		limit := 10
		offset := 0

		if p, err := strconv.Atoi(page); err == nil && p > 1 {
			offset = (p - 1) * limit
		}

		query := "SELECT id, name, description, deadline, created_at, updated_at FROM goals"
		var args []interface{}

		// Фильтрация по имени (ILIKE '%filter%')
		if filter != "" {
			query += " WHERE name ILIKE $1"
			args = append(args, "%"+filter+"%")
		}

		// Сортировка (разрешён только один из полей "name", "deadline", "created_at", "updated_at")
		if sortField != "" {
			allowedSorts := map[string]bool{
				"name":       true,
				"deadline":   true,
				"created_at": true,
				"updated_at": true,
			}
			if allowedSorts[strings.ToLower(sortField)] {
				query += " ORDER BY " + sortField
			} else {
				goalLog.WithField("sortField", sortField).Error("Invalid sort field")
				http.Error(w, "Invalid sort field", http.StatusBadRequest)
				return
			}
		}

		// Пагинация
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		args = append(args, limit, offset)

		rows, err := db.Query(query, args...)
		if err != nil {
			goalLog.WithFields(logrus.Fields{
				"error": err.Error(),
				"query": query,
				"args":  args,
			}).Error("Failed to retrieve goals")
			http.Error(w, "Failed to retrieve goals", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var goals []Goal
		for rows.Next() {
			var g Goal
			if err := rows.Scan(
				&g.ID,
				&g.Name,
				&g.Description,
				&g.Deadline,
				&g.CreatedAt,
				&g.UpdatedAt,
			); err != nil {
				goalLog.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Failed to scan goals row")
				http.Error(w, "Failed to scan goals", http.StatusInternalServerError)
				return
			}
			goals = append(goals, g)
		}

		// Если целей нет, вернём пустой массив []
		w.Header().Set("Content-Type", "application/json")
		if len(goals) == 0 {
			json.NewEncoder(w).Encode([]Goal{})
			return
		}

		// Иначе вернём цели
		json.NewEncoder(w).Encode(goals)
	}
}

// UpdateGoal — Обработчик для обновления цели (по старому имени)
func UpdateGoal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			OldName     string `json:"oldName"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Deadline    string `json:"deadline"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			goalLog.WithField("error", err.Error()).Error("Invalid input format for update")
			http.Error(w, "Invalid input format", http.StatusBadRequest)
			return
		}

		query := `
            UPDATE goals
            SET name = $1, description = $2, deadline = $3, updated_at = NOW()
            WHERE name = $4
        `
		res, err := db.Exec(query, input.Name, input.Description, input.Deadline, input.OldName)
		if err != nil {
			goalLog.WithFields(logrus.Fields{
				"error":   err.Error(),
				"oldName": input.OldName,
				"newName": input.Name,
			}).Error("Failed to update goal")
			http.Error(w, "Failed to update goal", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			goalLog.WithField("oldName", input.OldName).Warn("Goal with specified name not found")
			http.Error(w, "Goal with the specified name not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Goal successfully updated",
		})
	}
}
func DeleteGoalByName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		goalLog.Info("DeleteGoalByName called")

		// Парсинг JSON Body
		var input struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			goalLog.Warn("Invalid JSON format")
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		if input.Name == "" {
			goalLog.Warn("Missing 'name' field in request body")
			http.Error(w, "Goal name is required", http.StatusBadRequest)
			return
		}

		goalLog.Infof("Attempting to delete goal: %s", input.Name)

		query := `DELETE FROM goals WHERE name = $1`
		res, err := db.Exec(query, input.Name)
		if err != nil {
			goalLog.WithField("error", err.Error()).Error("Failed to delete goal")
			http.Error(w, "Failed to delete goal", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			goalLog.Infof("Goal with the name '%s' not found", input.Name)
			http.Error(w, "Goal with the specified name not found", http.StatusNotFound)
			return
		}

		goalLog.Infof("Goal '%s' deleted successfully", input.Name)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Goal successfully deleted",
		})
	}
}

// DeleteAllGoals — Обработчик для удаления всех целей
func DeleteAllGoals(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		goalLog.Info("DeleteAllGoals called")

		query := `DELETE FROM goals`
		res, err := db.Exec(query)
		if err != nil {
			goalLog.WithField("error", err.Error()).Error("Failed to delete all goals")
			http.Error(w, "Failed to delete all goals", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		goalLog.Infof("Deleted %d goals successfully", rowsAffected)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       "All goals successfully deleted",
			"rows_affected": rowsAffected,
		})
	}
}

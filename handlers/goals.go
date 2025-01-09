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
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&goal); err != nil {
			goalLog.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Invalid input for goal creation")
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		query := `INSERT INTO goals (name, description, deadline, created_at, updated_at) 
		          VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at`
		err := db.QueryRow(query, goal.Name, goal.Description, goal.Deadline).Scan(&goal.ID, &goal.CreatedAt, &goal.UpdatedAt)
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

// GetGoals — Обработчик для получения списка целей
func GetGoals(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("filter")
		sort := r.URL.Query().Get("sort")
		page := r.URL.Query().Get("page")

		limit := 10
		offset := 0

		if p, err := strconv.Atoi(page); err == nil && p > 1 {
			offset = (p - 1) * limit
		}

		query := "SELECT id, name, description, deadline, created_at, updated_at FROM goals"
		var args []interface{}

		if filter != "" {
			query += fmt.Sprintf(" WHERE name ILIKE $%d", len(args)+1)
			args = append(args, "%"+filter+"%")
		}

		if sort != "" {
			allowedSorts := map[string]bool{
				"name":       true,
				"deadline":   true,
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
			var goal Goal
			if err := rows.Scan(&goal.ID, &goal.Name, &goal.Description, &goal.Deadline, &goal.CreatedAt, &goal.UpdatedAt); err != nil {
				goalLog.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Failed to scan goals")
				http.Error(w, "Failed to scan goals", http.StatusInternalServerError)
				return
			}
			goals = append(goals, goal)
		}

		if len(goals) == 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]Goal{})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(goals)
	}
}

// UpdateGoal — Обработчик для обновления цели
func UpdateGoal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var goal struct {
			OldName     string `json:"oldName"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Deadline    string `json:"deadline"`
		}

		if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
			http.Error(w, "Invalid input format", http.StatusBadRequest)
			return
		}

		query := `UPDATE goals 
		          SET name = $1, description = $2, deadline = $3, updated_at = NOW() 
		          WHERE name = $4`
		res, err := db.Exec(query, goal.Name, goal.Description, goal.Deadline, goal.OldName)
		if err != nil {
			http.Error(w, "Failed to update goal", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Goal with the specified name not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Goal successfully updated",
		})
	}
}

// DeleteGoalByName — Обработчик для удаления цели по названию
func DeleteGoalByName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DeleteGoalByName called") // Отладка

		name := r.URL.Query().Get("name")
		if name == "" {
			log.Println("Missing 'name' parameter")
			http.Error(w, "Goal name is required", http.StatusBadRequest)
			return
		}

		log.Printf("Attempting to delete goal: %s", name) // Отладка

		query := `DELETE FROM goals WHERE name = $1`
		res, err := db.Exec(query, name)
		if err != nil {
			log.Printf("Failed to delete goal: %s", err.Error()) // Отладка
			http.Error(w, "Failed to delete goal", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			log.Printf("Goal with the name '%s' not found", name) // Отладка
			http.Error(w, "Goal with the specified name not found", http.StatusNotFound)
			return
		}

		log.Printf("Goal '%s' deleted successfully", name) // Отладка

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Goal successfully deleted",
		})
	}
}

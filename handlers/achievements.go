package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Achievement - структура для хранения данных о достижениях
type Achievement struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

// CreateAchievement - Обработчик для добавления достижения
func CreateAchievement(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var achievement Achievement
		if err := json.NewDecoder(r.Body).Decode(&achievement); err != nil {
			http.Error(w, "Invalid data format", http.StatusBadRequest)
			return
		}

		// SQL-запрос для добавления достижения
		query := `INSERT INTO achievements (user_id, title, description, date) 
		          VALUES ($1, $2, $3, $4) RETURNING id, user_id, title, description, date`
		err := db.QueryRow(query, achievement.UserID, achievement.Title, achievement.Description, achievement.Date).Scan(
			&achievement.ID, &achievement.UserID, &achievement.Title, &achievement.Description, &achievement.Date,
		)
		if err != nil {
			http.Error(w, "Error adding achievement: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(achievement)
	}
}

// GetAchievements - Обработчик для получения списка достижений
func GetAchievements(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, user_id, title, description, date FROM achievements")
		if err != nil {
			http.Error(w, "Error retrieving achievements", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var achievements []Achievement
		for rows.Next() {
			var achievement Achievement
			if err := rows.Scan(&achievement.ID, &achievement.UserID, &achievement.Title, &achievement.Description, &achievement.Date); err != nil {
				http.Error(w, "Error processing achievements", http.StatusInternalServerError)
				return
			}
			achievements = append(achievements, achievement)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(achievements)
	}
}

// DeleteAchievement - Обработчик для удаления достижения
func DeleteAchievementByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid achievement ID", http.StatusBadRequest)
			return
		}

		query := "DELETE FROM achievements WHERE id = $1"
		res, err := db.Exec(query, id)
		if err != nil {
			http.Error(w, "Error deleting achievement", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			http.Error(w, "Error checking deletion", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Achievement not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Achievement successfully deleted",
		})
	}
}

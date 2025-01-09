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

// Achievement — структура для достижения
type Achievement struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

var achievementLog = logrus.New()

func init() {
	achievementLog.SetFormatter(&logrus.JSONFormatter{})
	achievementLog.SetLevel(logrus.InfoLevel)
}

// CreateAchievement — Обработчик для добавления достижения
func CreateAchievement(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var achievement Achievement
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Disallow unknown fields

		if err := decoder.Decode(&achievement); err != nil {
			achievementLog.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Invalid input for achievement creation")
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		query := `INSERT INTO achievements (title, description, date) 
		          VALUES ($1, $2, NOW()) RETURNING id`
		err := db.QueryRow(query, achievement.Title, achievement.Description).Scan(&achievement.ID)
		if err != nil {
			achievementLog.WithFields(logrus.Fields{
				"error":       err.Error(),
				"title":       achievement.Title,
				"description": achievement.Description,
			}).Error("Failed to create achievement")
			http.Error(w, "Failed to create achievement", http.StatusInternalServerError)
			return
		}

		achievementLog.WithFields(logrus.Fields{
			"id":          achievement.ID,
			"title":       achievement.Title,
			"description": achievement.Description,
		}).Info("Achievement created successfully")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(achievement)
	}
}

// GetAchievements — Обработчик для получения списка достижений
func GetAchievements(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("filter") // Фильтрация по тексту
		sort := r.URL.Query().Get("sort")     // Сортировка по полю
		page := r.URL.Query().Get("page")     // Номер страницы

		limit := 10 // Количество записей на страницу
		offset := 0 // Сдвиг для SQL-запроса

		if p, err := strconv.Atoi(page); err == nil && p > 1 {
			offset = (p - 1) * limit
		}

		query := "SELECT id, title, description, date FROM achievements"
		var args []interface{}

		if filter != "" {
			query += fmt.Sprintf(" WHERE title ILIKE $%d", len(args)+1)
			args = append(args, "%"+filter+"%")
		}

		if sort != "" {
			allowedSorts := map[string]bool{
				"title":       true,
				"date":        true,
				"description": true,
			}
			if allowedSorts[strings.ToLower(sort)] {
				query += " ORDER BY " + sort
			} else {
				achievementLog.WithFields(logrus.Fields{
					"sort": sort,
				}).Error("Invalid sort field")
				http.Error(w, "Invalid sort field", http.StatusBadRequest)
				return
			}
		}

		argIndex := len(args) + 1
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, limit, offset)

		rows, err := db.Query(query, args...)
		if err != nil {
			achievementLog.WithFields(logrus.Fields{
				"error": err.Error(),
				"query": query,
				"args":  args,
			}).Error("Failed to retrieve achievements")
			http.Error(w, "Failed to retrieve achievements", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var achievements []Achievement
		for rows.Next() {
			var achievement Achievement
			if err := rows.Scan(&achievement.ID, &achievement.Title, &achievement.Description, &achievement.Date); err != nil {
				achievementLog.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Failed to scan achievements")
				http.Error(w, "Failed to scan achievements", http.StatusInternalServerError)
				return
			}
			achievements = append(achievements, achievement)
		}

		if len(achievements) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]Achievement{})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(achievements)
	}
}

// UpdateAchievement — Обработчик для обновления достижения
func UpdateAchievement(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Only PUT requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		var achievement struct {
			OldTitle    string `json:"oldTitle"`
			Title       string `json:"title"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&achievement); err != nil {
			achievementLog.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Invalid input format for updating achievement")
			http.Error(w, "Invalid input format", http.StatusBadRequest)
			return
		}

		query := `UPDATE achievements 
		          SET title = $1, description = $2, date = NOW() 
		          WHERE title = $3`

		achievementLog.WithFields(logrus.Fields{
			"newTitle":    achievement.Title,
			"description": achievement.Description,
			"oldTitle":    achievement.OldTitle,
		}).Info("Attempting to update achievement")

		res, err := db.Exec(query, achievement.Title, achievement.Description, achievement.OldTitle)
		if err != nil {
			achievementLog.WithFields(logrus.Fields{
				"error":    err.Error(),
				"oldTitle": achievement.OldTitle,
				"title":    achievement.Title,
			}).Error("Failed to update achievement")
			http.Error(w, "Failed to update achievement", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			achievementLog.WithFields(logrus.Fields{
				"oldTitle": achievement.OldTitle,
			}).Warn("Achievement with the specified title not found")
			http.Error(w, "Achievement with the specified title not found", http.StatusNotFound)
			return
		}

		achievementLog.WithFields(logrus.Fields{
			"oldTitle": achievement.OldTitle,
			"title":    achievement.Title,
		}).Info("Achievement updated successfully")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Achievement successfully updated",
		})
	}
}

// DeleteAchievementByTitle — Обработчик для удаления достижения
func DeleteAchievementByTitle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("title")
		if title == "" {
			http.Error(w, "Achievement title is required", http.StatusBadRequest)
			return
		}

		query := `DELETE FROM achievements WHERE title = $1`
		result, err := db.Exec(query, title)
		if err != nil {
			http.Error(w, "Failed to delete achievement", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Achievement not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Achievement deleted successfully"))
	}
}

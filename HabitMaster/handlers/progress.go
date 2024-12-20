package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// Progress - структура для записи прогресса
type Progress struct {
	ID      int    `json:"id"`
	HabitID int    `json:"habit_id"`
	Date    string `json:"date"`
	Status  string `json:"status"`
}

// CreateProgress - Обработчик для добавления записи прогресса
func CreateProgress(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Только POST-запросы разрешены", http.StatusMethodNotAllowed)
			return
		}

		var progress Progress
		if err := json.NewDecoder(r.Body).Decode(&progress); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		// SQL-запрос для добавления записи
		query := `INSERT INTO progress (habit_id, date, status) VALUES ($1, $2, $3)`
		_, err := db.Exec(query, progress.HabitID, progress.Date, progress.Status)
		if err != nil {
			http.Error(w, "Ошибка при добавлении записи прогресса", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Прогресс успешно добавлен"})
	}
}

// GetProgress - Обработчик для получения всех записей прогресса
func GetProgress(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Только GET-запросы разрешены", http.StatusMethodNotAllowed)
			return
		}

		// SQL-запрос для получения всех записей
		rows, err := db.Query(`SELECT id, habit_id, date, status FROM progress`)
		if err != nil {
			http.Error(w, "Ошибка при получении данных из базы", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var progressRecords []Progress
		for rows.Next() {
			var progress Progress
			if err := rows.Scan(&progress.ID, &progress.HabitID, &progress.Date, &progress.Status); err != nil {
				http.Error(w, "Ошибка при обработке данных из базы", http.StatusInternalServerError)
				return
			}
			progressRecords = append(progressRecords, progress)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(progressRecords)
	}
}

// UpdateProgressStatus - Обработчик для обновления статуса прогресса
func UpdateProgressStatus(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Только PUT-запросы разрешены", http.StatusMethodNotAllowed)
			return
		}

		var progress Progress
		if err := json.NewDecoder(r.Body).Decode(&progress); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		// SQL-запрос для обновления статуса
		query := `UPDATE progress SET status = $1 WHERE id = $2`
		res, err := db.Exec(query, progress.Status, progress.ID)
		if err != nil {
			http.Error(w, "Ошибка при обновлении статуса прогресса", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Запись не найдена", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Статус прогресса успешно обновлен"})
	}
}

package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Notification - структура для уведомлений
type Notification struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Message     string `json:"message"`
	ScheduledAt string `json:"scheduled_at"`
	IsSent      bool   `json:"is_sent"`
}

// CreateNotification - обработчик для добавления уведомления
func CreateNotification(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var notification Notification
		if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		if notification.UserID == 0 || notification.Message == "" || notification.ScheduledAt == "" {
			http.Error(w, "Обязательные поля: user_id, message, scheduled_at", http.StatusBadRequest)
			return
		}

		query := `INSERT INTO notifications (user_id, message, scheduled_at, is_sent) 
		          VALUES ($1, $2, $3, $4) RETURNING id`
		err := db.QueryRow(query, notification.UserID, notification.Message, notification.ScheduledAt, false).Scan(&notification.ID)
		if err != nil {
			http.Error(w, "Ошибка при добавлении уведомления", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notification)
	}
}

// GetNotifications - обработчик для получения списка уведомлений
func GetNotifications(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, user_id, message, scheduled_at, is_sent FROM notifications")
		if err != nil {
			http.Error(w, "Ошибка при получении данных из базы", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var notifications []Notification
		for rows.Next() {
			var notification Notification
			if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Message, &notification.ScheduledAt, &notification.IsSent); err != nil {
				http.Error(w, "Ошибка при обработке данных из базы", http.StatusInternalServerError)
				return
			}
			notifications = append(notifications, notification)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notifications)
	}
}

// UpdateNotificationStatus - обработчик для обновления статуса отправки
func DeleteNotificationByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID уведомления", http.StatusBadRequest)
			return
		}

		query := "DELETE FROM notifications WHERE id = $1"
		res, err := db.Exec(query, id)
		if err != nil {
			http.Error(w, "Ошибка при удалении уведомления", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			http.Error(w, "Ошибка при проверке удаления", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Уведомление не найдено", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Уведомление успешно удалено",
		})
	}
}

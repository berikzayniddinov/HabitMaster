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

// Notification — структура для уведомления
type Notification struct {
	ID          int    `json:"id"`
	Message     string `json:"message"`
	ScheduledAt string `json:"scheduled_at"`
	IsSent      bool   `json:"is_sent"`
}

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

// CreateNotification — Обработчик для создания уведомления
func CreateNotification(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var notification Notification
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Disallow unknown fields

		if err := decoder.Decode(&notification); err != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Invalid input for notification creation")
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		query := `INSERT INTO notifications (message, scheduled_at, is_sent) 
		          VALUES ($1, $2, false) RETURNING id`
		err := db.QueryRow(query, notification.Message, notification.ScheduledAt).Scan(&notification.ID)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error":    err.Error(),
				"message":  notification.Message,
				"schedule": notification.ScheduledAt,
			}).Error("Failed to create notification")
			http.Error(w, "Failed to create notification", http.StatusInternalServerError)
			return
		}

		log.WithFields(logrus.Fields{
			"id":       notification.ID,
			"message":  notification.Message,
			"schedule": notification.ScheduledAt,
		}).Info("Notification created successfully")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notification)
	}
}

// GetNotifications — Обработчик для получения списка уведомлений
func GetNotifications(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("filter") // Фильтрация по тексту
		sort := r.URL.Query().Get("sort")     // Сортировка по полю
		page := r.URL.Query().Get("page")     // Номер страницы

		limit := 10 // Количество записей на страницу
		offset := 0 // Сдвиг для SQL-запроса

		if p, err := strconv.Atoi(page); err == nil && p > 1 {
			offset = (p - 1) * limit
		}

		query := "SELECT id, message, scheduled_at, is_sent FROM notifications"
		var args []interface{}

		if filter != "" {
			query += fmt.Sprintf(" WHERE message ILIKE $%d", len(args)+1)
			args = append(args, "%"+filter+"%")
		}

		if sort != "" {
			allowedSorts := map[string]bool{
				"message":      true,
				"scheduled_at": true,
				"is_sent":      true,
			}
			if allowedSorts[strings.ToLower(sort)] {
				query += " ORDER BY " + sort
			} else {
				log.WithFields(logrus.Fields{
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
			log.WithFields(logrus.Fields{
				"error": err.Error(),
				"query": query,
				"args":  args,
			}).Error("Failed to retrieve notifications")
			http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var notifications []Notification
		for rows.Next() {
			var notification Notification
			if err := rows.Scan(&notification.ID, &notification.Message, &notification.ScheduledAt, &notification.IsSent); err != nil {
				log.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Failed to scan notifications")
				http.Error(w, "Failed to scan notifications", http.StatusInternalServerError)
				return
			}
			notifications = append(notifications, notification)
		}

		if len(notifications) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]Notification{})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(notifications)
	}
}

// UpdateNotification — Обработчик для обновления уведомления
func UpdateNotification(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Only PUT requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		var notification struct {
			OldMessage  string `json:"oldMessage"`
			Message     string `json:"message"`
			ScheduledAt string `json:"scheduled_at"`
			IsSent      bool   `json:"is_sent"`
		}
		if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Invalid input format for updating notification")
			http.Error(w, "Invalid input format", http.StatusBadRequest)
			return
		}

		query := `UPDATE notifications 
		          SET message = $1, scheduled_at = $2, is_sent = $3, updated_at = NOW() 
		          WHERE message = $4`

		log.WithFields(logrus.Fields{
			"newMessage":  notification.Message,
			"scheduledAt": notification.ScheduledAt,
			"isSent":      notification.IsSent,
			"oldMessage":  notification.OldMessage,
		}).Info("Attempting to update notification")

		res, err := db.Exec(query, notification.Message, notification.ScheduledAt, notification.IsSent, notification.OldMessage)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error":      err.Error(),
				"oldMessage": notification.OldMessage,
				"message":    notification.Message,
			}).Error("Failed to update notification")
			http.Error(w, "Failed to update notification", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			log.WithFields(logrus.Fields{
				"oldMessage": notification.OldMessage,
			}).Warn("Notification with the specified message not found")
			http.Error(w, "Notification with the specified message not found", http.StatusNotFound)
			return
		}

		log.WithFields(logrus.Fields{
			"oldMessage": notification.OldMessage,
			"message":    notification.Message,
		}).Info("Notification updated successfully")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Notification successfully updated",
		})
	}
}

// DeleteNotificationByName — Обработчик для удаления уведомления по имени
// DeleteNotificationByName — Обработчик для удаления уведомления по имени
func DeleteNotificationByName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("message")
		if name == "" {
			http.Error(w, "Message parameter is required", http.StatusBadRequest)
			return
		}

		query := `DELETE FROM notifications WHERE message = $1`
		result, err := db.Exec(query, name)
		if err != nil {
			http.Error(w, "Failed to delete notification", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Notification deleted successfully"))
	}
}

package handlers

import (
	"HabitMaster/emailSender"
	"encoding/json"
	"net/http"
)

// MassEmailRequest структура для массовой рассылки
type MassEmailRequest struct {
	Emails  []string `json:"emails"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

// SupportEmailRequest структура для отправки в поддержку
type SupportEmailRequest struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// SendMassEmailHandler обработчик для массовой рассылки
func SendMassEmailHandler(emailSender *emailSender.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MassEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := emailSender.SendEmail(req.Emails, req.Subject, req.Body); err != nil {
			http.Error(w, "Failed to send emails", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Emails sent successfully"))
	}
}

// SendSupportEmailHandler обработчик для отправки письма в поддержку
func SendSupportEmailHandler(emailSender *emailSender.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SupportEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		adminEmail := "admin@example.com" // Замените на email администратора
		if err := emailSender.SendEmail([]string{adminEmail}, req.Subject, req.Body); err != nil {
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Support email sent successfully"))
	}
}

package handlers

import (
	"HabitMaster/emailSender"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
		r.ParseMultipartForm(10 << 20) // Ограничение размера файла (10MB)

		// Получаем данные формы
		recipients := r.FormValue("recipients")
		subject := r.FormValue("subject")
		body := r.FormValue("body")

		// Преобразуем строку получателей в массив
		to := strings.Split(recipients, ",")

		// Получаем файл
		file, handler, err := r.FormFile("attachment")
		var fileBytes []byte
		var fileName string
		if err == nil {
			defer file.Close()
			fileBytes, err = ioutil.ReadAll(file)
			if err != nil {
				http.Error(w, "Failed to read file", http.StatusInternalServerError)
				return
			}
			fileName = handler.Filename
		}

		// Отправляем email с файлом (если он был)
		log.Printf("Sending email to: %v, subject: %s, file: %s, file size: %d bytes", to, subject, fileName, len(fileBytes))
		err = emailSender.SendEmailWithAttachment(
			to,        // Список получателей
			subject,   // Тема письма
			body,      // Тело письма
			fileName,  // Имя файла
			fileBytes, // Содержимое файла
		)
		if err != nil {
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Emails sent successfully")
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

// SendEmailWithAttachmentHandler обработчик для отправки email с вложением
func SendEmailWithAttachmentHandler(emailSender *emailSender.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20) // Ограничение на размер данных 10 MB
		if err != nil {
			http.Error(w, "Unable to parse form data", http.StatusBadRequest)
			return
		}

		// Получаем значения из формы
		recipients := r.FormValue("recipients")
		subject := r.FormValue("subject")
		body := r.FormValue("body")

		// Получение файла
		file, header, err := r.FormFile("attachment")
		if err != nil {
			http.Error(w, "Failed to get file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Чтение содержимого файла
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		// Отправка email с вложением
		if err := emailSender.SendEmailWithAttachment(
			strings.Split(recipients, ","), // Разделяем строку с email'ами на массив
			subject,
			body,
			header.Filename, // Имя файла
			fileBytes,       // Содержимое файла
		); err != nil {
			http.Error(w, "Failed to send email with attachment", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Email with attachment sent successfully"))
	}
}

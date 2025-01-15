package emailSender

import (
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"os"
	"strconv"
)

// EmailSender содержит параметры SMTP
type EmailSender struct {
	Host     string
	Port     int
	Username string
	Password string
}

// NewEmailSender создает новый экземпляр EmailSender
func NewEmailSender() *EmailSender {
	return &EmailSender{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     mustParseInt(os.Getenv("SMTP_PORT")),
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}

// mustParseInt преобразует строку в int, с обработкой ошибки
func mustParseInt(s string) int {
	port, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Invalid port value: %v", err)
	}
	return port
}

// SendEmail отправляет email
func (e *EmailSender) SendEmail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.Username)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// SendEmailWithAttachment отправляет email с вложением
func (e *EmailSender) SendEmailWithAttachment(to []string, subject, body, fileName string, fileData []byte) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.Username)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if len(fileData) > 0 && fileName != "" {
		m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(fileData)
			return err
		}))
	}

	d := gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err) // Лог ошибки
		return err
	}

	log.Printf("Email successfully sent to: %v", to) // Лог успешной отправки
	return nil
}

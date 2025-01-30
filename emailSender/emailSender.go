package emailSender

import (
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"os"
	"strconv"
)

type EmailSender interface {
	SendEmail(to []string, subject, body string) error
	SendEmailWithAttachment(to []string, subject, body, fileName string, fileData []byte) error
}

type RealEmailSender struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailSender() EmailSender {
	return &RealEmailSender{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     mustParseInt(os.Getenv("SMTP_PORT")),
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}

func mustParseInt(s string) int {
	port, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Invalid port value: %v", err)
	}
	return port
}

func (e *RealEmailSender) SendEmail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.Username)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)

	return d.DialAndSend(m)
}

func (e *RealEmailSender) SendEmailWithAttachment(to []string, subject, body, fileName string, fileData []byte) error {
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

	return d.DialAndSend(m)
}

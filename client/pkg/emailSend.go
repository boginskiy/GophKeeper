package pkg

import (
	"fmt"
	"net/smtp"
)

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type EmailSend struct {
	SMTPHost    string
	SMTPPort    string
	EmailFrom   string
	AppPassword string
	auth        smtp.Auth
}

func NewEmailSend(host, port, from, password string) *EmailSend {
	return &EmailSend{
		SMTPHost:    host,
		SMTPPort:    port,
		EmailFrom:   from,
		AppPassword: password,

		// Аутентификация
		auth: smtp.PlainAuth("", from, password, host),
	}
}

func (e *EmailSend) SendEmail(to, subject, body string) error {
	// Формируем сообщение
	message := e.UseSimpleFormat(to, subject, body)

	// Отправка
	err := smtp.SendMail(e.SMTPHost+":"+e.SMTPPort, e.auth, e.EmailFrom, []string{to}, message)
	return err
}

func (e *EmailSend) UseSimpleFormat(to, subject, body string) []byte {
	where := fmt.Sprintf("To: %s\r\n", to)
	topic := fmt.Sprintf("Subject: %s\r\n", subject)
	contentType := "Content-Type: text/plain; charset=UTF-8\r\n"
	contentData := fmt.Sprintf("\r\n%s\r\n", body)

	return []byte(fmt.Sprint(
		where,
		topic,
		contentType,
		contentData))
}

package services

import (
	"bank-system/config"
	"fmt"
	"log"
	"net/smtp"
)

func SendEmailNotification(to, subject, body string) error {
	config := config.LoadSMTPConfig()

	if config.Host == "smtp.example.com" {
		log.Printf("SMTP not configured. Skipping email to %s: Subject: %s", to, subject)
		return nil
	}

	auth := smtp.PlainAuth(
		"",
		config.User,
		config.Password,
		config.Host,
	)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		config.From,
		to,
		subject,
		body,
	)

	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	err := smtp.SendMail(
		addr,
		auth,
		config.From,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		log.Printf("Error sending email to %s: %v", to, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", to)

	return nil
}

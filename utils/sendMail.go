package utils

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendMail(email, subject, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_USERNAME"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	// Set up the email server connection information
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"))

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	} else {
		return nil
	}
}

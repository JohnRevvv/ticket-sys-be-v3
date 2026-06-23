package email

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendOTP(to string, otp string) error {

	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := "587"

	if from == "" || password == "" || smtpHost == "" {
		return fmt.Errorf("email config missing in environment")
	}

	subject := "Your OTP Code"

	body := fmt.Sprintf(
		"Your OTP code is: %s\n\nThis code will expire in 5 minutes.",
		otp,
	)

	message := []byte(
		"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{to},
		message,
	)

	if err != nil {
		return err
	}

	return nil
}
package client

// sendmail client

import (
	"moodle/config"
	"net/smtp"
)

type SendmailClient struct {
}

// new sendmail client
func NewSendmailClient() *SendmailClient {
	return &SendmailClient{}
}


func (client *SendmailClient) SendMail(to []string, subject, body string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		config.EMAIL,
		config.EMAIL_PASSWORD,
		config.SMTP,
	)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		config.SMTP+":"+config.SMTP_PORT,
		auth,
		config.EMAIL,
		to,
		[]byte(body),
	)
	if err != nil {
		return err
	}

	return nil
}

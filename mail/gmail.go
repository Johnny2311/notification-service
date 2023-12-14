package mail

import (
	"fmt"
	"net/smtp"
)

// same implementation for both Google and Yahoo just to remark Strategy pattern
// in the practice, just changing the config (envs) we can use the same code

type googleMailProvider struct {
	address    string
	password   string
	smtpServer string
	smtpPort   int64
}

func NewGoogleMailProvider() Mailer {
	return &googleMailProvider{}
}

func (g googleMailProvider) Send(to string, subject string, body string) error {

	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	auth := smtp.PlainAuth("", g.address, g.password, g.smtpServer)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", g.smtpServer, g.smtpPort),
		auth,
		g.address,
		[]string{to},
		[]byte(message),
	)
	if err != nil {
		fmt.Println("Error sending email:", err)
	}

	return nil
}

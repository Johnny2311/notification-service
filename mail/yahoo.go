package mail

import (
	"fmt"
	"net/smtp"
)

// same implementation for both Google and Yahoo just to remark Strategy pattern
// in the practice, just changing the config (envs) we can use the same code

type yahooMailProvider struct {
	address    string
	password   string
	smtpServer string
	smtpPort   int64
}

func NewYahooMailProvider() Mailer {
	return &yahooMailProvider{}
}

func (y yahooMailProvider) Send(to string, subject string, body string) error {

	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	auth := smtp.PlainAuth("", y.address, y.password, y.smtpServer)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", y.smtpServer, y.smtpPort),
		auth,
		y.address,
		[]string{to},
		[]byte(message),
	)
	if err != nil {
		fmt.Println("Error sending email:", err)
	}

	return nil
}

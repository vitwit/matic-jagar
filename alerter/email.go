package alerter

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmail to send mail alert
func (e emailAlert) SendEmail(msg, token, toEmail, fromEmail, accountName string) error {
	from := mail.NewEmail(accountName, fromEmail) //mail.NewEmail("Matic Tool", "matic@vitwit.com")
	subject := msg
	to := mail.NewEmail(accountName, toEmail) //mail.NewEmail("Matic Tool", toEmail)
	plainTextContent := msg
	htmlContent := msg
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(token)
	_, err := client.Send(message)
	if err != nil {
		return err
	}
	return nil
}

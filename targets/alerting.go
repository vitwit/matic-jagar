package targets

import (
	"log"
	"strconv"
	"strings"

	"github.com/vitwit/matic-jagar/alerter"
	"github.com/vitwit/matic-jagar/config"
)

// SendTelegramAlert sends the alert to telegram account
// check's alert setting before sending the alert
func SendTelegramAlert(msg string, cfg *config.Config) error {
	if strings.ToUpper(strconv.FormatBool(cfg.EnableAlerts.EnableTelegramAlerts)) == "TRUE" {
		if err := alerter.NewTelegramAlerter().SendTelegramMessage(msg, cfg.Telegram.BotToken, cfg.Telegram.ChatID); err != nil {
			log.Printf("failed to send tg alert: %v", err)
			return err
		}
	}
	return nil
}

// SendEmailAlert sends alert to email account
//by checking user's choice
func SendEmailAlert(msg string, cfg *config.Config) error {
	if strings.ToUpper(strconv.FormatBool(cfg.EnableAlerts.EnableEmailAlerts)) == "TRUE" {
		fromMail := cfg.SendGrid.SendgridEmail
		accountName := cfg.SendGrid.SendgridName
		if err := alerter.NewEmailAlerter().SendEmail(msg, cfg.SendGrid.Token, cfg.SendGrid.ToEmailAddress, fromMail, accountName); err != nil {
			log.Printf("failed to send email alert: %v", err)
			return err
		}
	}
	return nil
}

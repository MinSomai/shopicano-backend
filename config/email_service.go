package config

import (
	"github.com/spf13/viper"
)

type EmailServiceCfg struct {
	SMTPHost         string
	SMTPPort         int
	SMTPUsername     string
	SMTPPassword     string
	FromEmailAddress string
}

var emailServiceCfg EmailServiceCfg

func LoadEmailService() {
	mu.Lock()
	defer mu.Unlock()

	emailServiceCfg = EmailServiceCfg{
		SMTPHost:         viper.GetString("email_service.smtp_host"),
		SMTPPort:         viper.GetInt("email_service.smtp_port"),
		SMTPUsername:     viper.GetString("email_service.smtp_username"),
		SMTPPassword:     viper.GetString("email_service.smtp_password"),
		FromEmailAddress: viper.GetString("email_service.from_email_address"),
	}
}

func EmailService() EmailServiceCfg {
	return emailServiceCfg
}

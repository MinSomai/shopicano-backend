package services

import (
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/shopicano/shopicano-backend/config"
)

func SendSignUpVerificationEmail(name, email, userID, verificationToken string) error {
	body := fmt.Sprintf(`Hello %s,<br><br>`+
		`Please <a href="%s/v1/email-verification?uid=%s&token=%s">click here to verify your email.</a><br><br><br>`+
		`With love,<br/>`+
		`Shopicano Ltd.`, name, config.EmailService().VerificationUrl, userID, verificationToken)

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s", config.EmailService().FromEmailAddress))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Shopicano Verification Email")
	m.SetBody("text/html", body)

	if err := EmailDialer().DialAndSend(m); err != nil {
		return err
	}
	return nil
}

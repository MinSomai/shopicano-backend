package services

import (
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/matcornic/hermes/v2"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/config"
)

func SendSignUpVerificationEmail(name, email, userID, verificationToken string) error {
	emailTemplate := hermes.Email{
		Body: hermes.Body{
			Greeting:  "Hello",
			Signature: "With Love",
			Name:      name,
			Actions: []hermes.Action{
				{
					Instructions: "Thanks for using our service. Please,",
					Button: hermes.Button{
						Link: fmt.Sprintf("%s/v1/email-verification?uid=%s&token=%s",
							config.EmailService().VerificationUrl, userID, verificationToken),
						Color:     "green",
						Text:      "Activate your account",
						TextColor: "black",
					},
				},
			},
		},
	}

	body, err := app.Hermes().GenerateHTML(emailTemplate)
	if err != nil {
		return err
	}

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

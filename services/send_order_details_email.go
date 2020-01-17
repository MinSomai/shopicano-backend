package services

import (
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/matcornic/hermes/v2"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/config"
	"github.com/shopicano/shopicano-backend/models"
)

func SendOrderDetailsEmail(name, email string, order *models.OrderDetailsView) error {
	var entries [][]hermes.Entry

	for _, v := range order.Items {
		var rows []hermes.Entry
		rows = append(rows, hermes.Entry{
			Key:   "Item",
			Value: v.Name,
		})
		rows = append(rows, hermes.Entry{
			Key:   "Quantity",
			Value: fmt.Sprintf("%d", v.Quantity),
		})
		rows = append(rows, hermes.Entry{
			Key:   "Price",
			Value: fmt.Sprintf("%d", v.Price),
		})
		rows = append(rows, hermes.Entry{
			Key:   "Sub Total",
			Value: fmt.Sprintf("%d", v.SubTotal),
		})

		entries = append(entries, rows)
	}

	// Tax
	var taxRows []hermes.Entry
	taxRows = append(taxRows, hermes.Entry{
		Key:   "Item",
		Value: "",
	})
	taxRows = append(taxRows, hermes.Entry{
		Key:   "Quantity",
		Value: "",
	})
	taxRows = append(taxRows, hermes.Entry{
		Key:   "Price",
		Value: "Total Additional Charge",
	})
	taxRows = append(taxRows, hermes.Entry{
		Key:   "Sub Total",
		Value: fmt.Sprintf("%d", order.TotalAdditionalCharge),
	})
	entries = append(entries, taxRows)

	// Grand Total
	var rows []hermes.Entry
	rows = append(rows, hermes.Entry{
		Key:   "Item",
		Value: "",
	})
	rows = append(rows, hermes.Entry{
		Key:   "Quantity",
		Value: "",
	})
	rows = append(rows, hermes.Entry{
		Key:   "Price",
		Value: "Grand Total",
	})
	rows = append(rows, hermes.Entry{
		Key:   "Sub Total",
		Value: fmt.Sprintf("%d", order.GrandTotal),
	})
	entries = append(entries, rows)

	emailTemplate := hermes.Email{
		Body: hermes.Body{
			Greeting:  "Hello",
			Signature: "With Love",
			Name:      name,
			Intros: []string{
				"Your order has been placed successfully.",
			},
			Table: hermes.Table{
				Data: entries,
			},
			Actions: []hermes.Action{
				{
					Instructions: "",
					Button: hermes.Button{
						Link:      fmt.Sprintf("%s/v1/orders/%s", config.EmailService().VerificationUrl, order.ID),
						Color:     "white",
						Text:      "Order Details",
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
	m.SetHeader("Subject", "Order has been placed")
	m.SetBody("text/html", body)

	if err := EmailDialer().DialAndSend(m); err != nil {
		return err
	}
	return nil
}

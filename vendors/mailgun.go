package vendors

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/granitebps/crypto-price-alert/types"
	"github.com/mailgun/mailgun-go/v4"
)

func SendEmail(alert types.Alert, price int) error {
	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAILGUN_API_KEY")
	sender := os.Getenv("MAILGUN_SENDER")

	if domain != "" && apiKey != "" {
		mg := mailgun.NewMailgun(domain, apiKey)
		subject := fmt.Sprintf("Price Alert for %s!!!", alert.Ticker)
		body := fmt.Sprintf("%s price now is : %d", alert.Ticker, price)

		message := mg.NewMessage(sender, subject, body, alert.Email)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		// Send the message with a 10 second timeout
		resp, id, err := mg.Send(ctx, message)

		if err != nil {
			return err
		}

		fmt.Printf("Send to email. ID: %s Resp: %s\n", id, resp)
	}
	return nil
}

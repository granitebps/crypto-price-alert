package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/granitebps/crypto-price-alert/helper"
	"github.com/granitebps/crypto-price-alert/types"
	"github.com/granitebps/crypto-price-alert/vendors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}
	vendor := os.Getenv("PRICE_VENDOR")
	mail := os.Getenv("MAIL_VENDOR")
	sentryDsn := os.Getenv("SENTRY_DSN")
	appEnv := os.Getenv("APP_ENV")
	appDebug, _ := strconv.ParseBool(os.Getenv("APP_DEBUG"))

	err = sentry.Init(sentry.ClientOptions{
		Dsn:         sentryDsn,
		Debug:       appDebug,
		Environment: appEnv,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	app := fiber.New()
	app.Get("", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Success",
		})
	})

	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Minute().Do(func() {
		// Healthcheck
		apiUrl := os.Getenv("HEALTHCHECK_URL")
		request, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			log.Println(err)
			sentry.CaptureException(err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)

		if err != nil {
			log.Println(err)
			sentry.CaptureException(err)
			return
		}

		defer response.Body.Close()
	})

	s.Every(1).Hour().Do(func() {
		var alerts []types.Alert
		var lastPrice int

		filename := "alert.json"
		jsonData, err := helper.ReadJsonFile(filename)
		if err != nil {
			log.Println(err)
			sentry.CaptureException(err)
			return
		}

		result, err := helper.ParseAlertData(jsonData, alerts)
		if err != nil {
			log.Println(err)
			sentry.CaptureException(err)
			return
		}

		for _, alert := range result {
			if vendor == "indodax" {
				lastPrice, err = vendors.GetPriceIndodax(alert)
				if err != nil {
					log.Println(err)
					sentry.CaptureException(err)
					continue
				}
			} else if vendor == "coingecko" {
				lastPrice, err = vendors.GetPriceCoingecko(alert)
				if err != nil {
					log.Println(err)
					sentry.CaptureException(err)
					continue
				}
			} else {
				err := errors.New("price vendor not available")
				log.Println(err)
				sentry.CaptureException(err)
				continue
			}

			now := time.Now().Format(time.RFC1123)
			fmt.Printf("Price for %s at %s is : %d\n", alert.Pair, now, lastPrice)

			if lastPrice >= int(alert.Price) {
				if alert.Enabled {
					if mail == "mailgun" {
						err := vendors.SendEmail(alert, lastPrice)
						if err != nil {
							log.Println(err)
							sentry.CaptureException(err)
							continue
						}
					} else {
						err := errors.New("mail vendor not available")
						log.Println(err)
						sentry.CaptureException(err)
						continue
					}
				}
			}
		}
	})

	s.StartAsync()

	log.Fatal(app.Listen(":8000"))
}

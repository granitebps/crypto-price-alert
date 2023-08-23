package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
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
	defer sentry.Flush(2 * time.Second)
	defer sentry.Recover()

	var alerts []types.Alert
	var lastPrice int

	filename := "alert.json"
	jsonData, err := helper.ReadJsonFile(filename)
	if err != nil {
		log.Panic(err)
	}

	result, err := helper.ParseAlertData(jsonData, alerts)
	if err != nil {
		log.Panic(err)
	}

	for _, alert := range result {
		if vendor == "indodax" {
			lastPrice, err = vendors.GetPriceIndodax(alert)
			if err != nil {
				log.Panic(err)
			}
		} else if vendor == "coingecko" {
			lastPrice, err = vendors.GetPriceCoingecko(alert)
			if err != nil {
				log.Panic(err)
			}
		} else {
			log.Panic("Price vendor not available!")
		}

		now := time.Now().Format(time.RFC1123)
		fmt.Printf("Price for %s at %s is : %d\n", alert.Pair, now, lastPrice)

		if lastPrice >= int(alert.Price) {
			if alert.Enabled {
				if mail == "mailgun" {
					err := vendors.SendEmail(alert, lastPrice)
					if err != nil {
						log.Panic(err)
					}
				} else {
					log.Panic("Mail vendor not available!")
				}
			}
		}
	}
}

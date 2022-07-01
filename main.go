package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	var alerts []types.Alert
	var lastPrice int

	filename := "alert.json"
	jsonData, err := helper.ReadJsonFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	result, err := helper.ParseAlertData(jsonData, alerts)
	if err != nil {
		log.Fatal(err)
	}

	for i, alert := range result {
		if vendor == "indodax" {
			lastPrice, err = vendors.GetPrice(alert)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Price vendor not available!")
		}

		now := time.Now().Format(time.RFC1123)
		fmt.Printf("Price for %s at %s is : %d\n", alert.Pair, now, lastPrice)

		if lastPrice >= int(alert.Price) {
			if !helper.DateEqual(time.Now(), alert.EmailedAt) || alert.EmailedAt.IsZero() {
				if mail == "mailgun" {
					err := vendors.SendEmail(alert, lastPrice)
					if err != nil {
						log.Fatal(err)
					}
				} else {
					log.Fatal("Mail vendor not available!")
				}
			}
		}

		result[i].EmailedAt = time.Now()
	}
	helper.UpdateFile(result)
}

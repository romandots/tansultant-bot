package main

import (
	"log"
	"os"
	"tansulbot/internal/models"
)

var app *models.App

func main() {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Fatalf("Set the telegram token first")
	}

	var err error
	app, err = models.InitApp(telegramToken)
	if err != nil {
		log.Fatalf("Error initializing app: %v", err)
	}
	defer app.Close()

	go app.WaitForErrors()
	app.ReadTelegramUpdates()
}

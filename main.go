package main

import (
	"log"
	"os"
	"tansulbot/pkg/app"
)

var a *app.App

func main() {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Fatalf("Set the telegram token first")
	}

	var err error
	a, err = app.InitApp(telegramToken)
	if err != nil {
		log.Fatalf("Error initializing app: %v", err)
	}
	defer a.Close()

	// Drop all data from the database
	//err = a.DropAll()
	//if err != nil {
	//	fmt.Println("Error dropping database:", err)
	//	return
	//}

	go a.WaitForErrors()
	a.ReadTelegramUpdates()
}

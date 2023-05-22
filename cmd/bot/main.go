package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"tansulbot/pkg/bot"
)

func main() {
	fmt.Println("Starting work")

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Println("Set the telegram token first")
		return
	}

	c := bot.NewClient(telegramToken)

	go func() {
		for commandError := range c.Errors {
			if commandError != nil {
				fmt.Println("Error running command:", commandError.Error())
			}
		}
	}()

	for {
		updates, err := c.GetUpdates()
		if err != nil {
			fmt.Println("Error getting updates:", err)
		}

		c.ReadUpdates(updates.Result)

		time.Sleep(1 * time.Second)
	}
}

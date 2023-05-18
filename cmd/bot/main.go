package main

import (
	"fmt"
	"log"

	"tansulbot/pkg/bot"
)

const (
	telegramToken = "SECRET_TOKEN"
)

func main() {
	c := bot.NewClient(telegramToken)
	updates, err := c.GetUpdates()
	if err != nil {
		log.Println("Error getting updates:", err)
		return
	}

	fmt.Println(updates)
}

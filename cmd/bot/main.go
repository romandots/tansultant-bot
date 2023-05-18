package bot

import (
	"fmt"
	"log"
)

func main() {
	c := NewClient(token)
	updates, err := c.GetUpdates()
	if err != nil {
		log.Println("Error getting updates:", err)
		return
	}

	fmt.Println(updates)
}

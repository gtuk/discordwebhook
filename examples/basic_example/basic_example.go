package main

import (
	"log"

	"github.com/gtuk/discordwebhook"
)

func main() {
	var username = "BotUser"
	var content = "This is a test message"
	var url = ""

	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}
	for {
		err := discordwebhook.SendMessageRateLimitAware(url, message)
		if err != nil {
			log.Fatal(err)
		}

	}
}

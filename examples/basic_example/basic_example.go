package main

import (
	"log"

	"github.com/gtuk/discordwebhook"
)

func main() {
	var username = "BotUser"
	var content = "This is a test message"
	var url = "https://discord.com/api/webhooks/..."
	r1 := discordwebhook.NewRatelimiter()

	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}

	err := discordwebhook.SendMessage(url, message, r1)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/gtuk/discordwebhook"
)

func main() {
	var username = "BotUser"
	var content = "This is a test message"
	var url = "https://discord.com/api/webhooks/..."
	var filePath = "/path/to/file"

	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}

	err := discordwebhook.SendMessage(url, message)
	if err != nil {
		log.Fatal(err)
	}

	err = discordwebhook.SendFile(url, filePath)
	if err != nil {
		log.Fatal(err)
	}
}


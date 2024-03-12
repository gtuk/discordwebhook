package main

import (
	"log"

	"github.com/gtuk/discordwebhook"
)

func main() {
	var username = "BotUser"
	var content = "This is a test message"
	var url = "https://discord.com/api/webhooks/1178791404290445363/af5LrdlRfzw_h80IEA7LAUfR6xvP3xGzD-Tw-jHiDmNuO7geESWjBIgyrC-pzVCJSxmA"

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

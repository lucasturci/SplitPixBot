package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	function "github.com/lucasturci/SplitPixBot"
)

func disableWebhook(telegramToken string) {
	log.Print("Disabling webhook!")
	disableWebhookurl := fmt.Sprintf("https://api.telegram.org/bot%s/deleteWebhook", telegramToken)
	http.Get(disableWebhookurl)
}

func enableWebhook(telegramToken, url string) {
	log.Print("Reenabling webhook!")
	enableWebhookurl := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=%s", telegramToken, url)
	http.Get(enableWebhookurl)
}

func onExit(telegramToken string) {
	content, err := ioutil.ReadFile("../webhook-url.txt")
	if err != nil {
		log.Printf("Something went wrong when reading webhook url local file: %s", err.Error())
	}
	webhookUrl := string(content)
	enableWebhook(telegramToken, webhookUrl)
}

func main() {
	// Read telegram token
	var telegramToken string
	if os.Getenv("TELEGRAM_TOKEN") != "" {
		telegramToken = os.Getenv("TELEGRAM_TOKEN")
	} else {
		log.Printf("Telegram Token environment variable not set, defaulting to reading from telegram-token.txt")
		content, err := ioutil.ReadFile("../telegram-token.txt")
		if err != nil {
			log.Printf("Something went wrong when reading telegram token local file: %s", err.Error())
		}
		telegramToken = string(content)
	}
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Printf("Error with telegram authentication: %s", err.Error())
		return
	}

	bot.Debug = true

	disableWebhook(telegramToken)

	// Gracefully reenable webhook on exit
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(telegramToken string) {
		<-c
		onExit(telegramToken)
		os.Exit(1)
	}(telegramToken)

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		function.HandleUpdate(bot, &update)
	}

}

package commands

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var CommandHandlerMap map[string]func(*tgbotapi.BotAPI, *tgbotapi.Message)

func startCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Vamos lá, o que você quer fazer?")

	var buttons = [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("Cadastrar um Pix", "/start-1")},
		{tgbotapi.NewInlineKeyboardButtonData("Cobrar alguém", "/start-2")},
		{tgbotapi.NewInlineKeyboardButtonData("Dividir contas de um grupo", "/start-3")},
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error when sending message: %s", err.Error())
	}
}

func helpCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

}

func init() {
	CommandHandlerMap = map[string]func(*tgbotapi.BotAPI, *tgbotapi.Message){
		"/start": startCommand,
		"/help":  helpCommand,
	}
}

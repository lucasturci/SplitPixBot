package commands

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var CommandHandlerMap map[string]func(*tgbotapi.BotAPI, *tgbotapi.Message)
var CallbackQueryMap map[string]func(*tgbotapi.BotAPI, *tgbotapi.CallbackQuery)

func startCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Vamos lá, o que você quer fazer?")

	var buttons = [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("Cadastrar um Pix", "start-1")},
		{tgbotapi.NewInlineKeyboardButtonData("Cobrar alguém", "start-2")},
		{tgbotapi.NewInlineKeyboardButtonData("Dividir contas de um grupo", "start-3")},
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error when sending message: %s", err.Error())
	}
}

func helpCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

}

func myPixCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Beleza, vamos configurar o pix")
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error when sending message: %s", err.Error())
	}
}

func myPixQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	myPixCommand(bot, query.Message)
	bot.Send(tgbotapi.NewCallback(query.ID, ""))
}

func chargeInstructionsQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, `Para cobrar alguém, 
	mande uma mensagem para a pessoa no seguinte formato 
	“@SplitPixBot <título> <valor>” (sem as aspas), onde <título> é um nome 
	para a despesa, e <valor> é o valor que você está cobrando, em reais. 
	Exemplo: “@SplitPixBot conta de gás 121,90”`)

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error when sending message: %s", err.Error())
	}
	bot.Send(tgbotapi.NewCallback(query.ID, ""))
}

func splitGroupBillsInstructionsQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, `Para dividar contas 
	entre membros de um grupo, é só me adicionar no grupo e iniciar uma 
	caderneta com /novacaderneta@SplitPixBot, onde você poderá adicionar 
	os membros. Depois, registre uma conta ou despesa 
	com /novadespesa@SplitPixBot. `)

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error when sending message: %s", err.Error())
	}
	bot.Send(tgbotapi.NewCallback(query.ID, ""))
}

func init() {
	CommandHandlerMap = map[string]func(*tgbotapi.BotAPI, *tgbotapi.Message){
		"start":  startCommand,
		"help":   helpCommand,
		"meupix": myPixCommand,
	}

	CallbackQueryMap = map[string]func(*tgbotapi.BotAPI, *tgbotapi.CallbackQuery){
		"start-1": myPixQuery,
		"start-2": chargeInstructionsQuery,
		"start-3": splitGroupBillsInstructionsQuery,
	}
}

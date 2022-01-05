package function

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lucasturci/SplitPixBot/commands"
)

func echo(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	log.Printf("Replying to message ID (%v) with body (%v)", update.Message.MessageID, update.Message.Text)

	// Now that we know we've gotten a new message, we can construct a
	// reply! We'll take the Chat ID and Text from the incoming message
	// and use it to create a new message.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	// We'll also say that this message is a reply to the previous message.
	// For any other specifications than Chat ID or Text, you'll need to
	// set fields on the `MessageConfig`.
	msg.ReplyToMessageID = update.Message.MessageID

	// Okay, we're sending our message off! We don't care about the message
	// we just sent, so we'll discard it.
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error when sending message reply to message id (%v): %s", update.Message.MessageID, err.Error())
		return err
	}
	return nil
}

func HandleUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {

	// We will only ever react to message updates
	if update.Message == nil {
		return nil
	}

	log.Printf("AAAAAAAA %+v", update.Message)

	// If it is a command, then let's parse it. Else, we will echo the message
	if update.Message.IsCommand() {
		command := update.Message.Command()
		commandHandler, exists := commands.CommandHandlerMap[command]
		if exists {
			commandHandler(bot, update.Message)
		} else {
			log.Printf("Unrecognized command %s. Ignoring", command)
		}
		return nil
	}

	return echo(bot, update)
}

func Handle(w http.ResponseWriter, r *http.Request) {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Printf("Error with telegram authentication: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bot.Debug = true

	update, err := bot.HandleUpdate(r)
	if err != nil {
		log.Printf("Error when parsing update: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := HandleUpdate(bot, update); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error when handling update")
	}

}

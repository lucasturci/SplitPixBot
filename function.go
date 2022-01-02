package function

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

	if update.Message == nil {
		return
	}

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
		// Note that panics are a bad way to handle errors. Telegram can
		// have service outages or network errors, you should retry sending
		// messages or more gracefully handle failures.
		log.Printf("Error with sending message reply to message id (%v): %s", update.Message.MessageID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hbourgeot/quoteme/tgbot"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	bot, err := tgbot.NewBotAPI(os.Getenv("TELEGRAM_BOT_URL"))
	if err != nil {

		fmt.Println("h", os.Getenv("TELEGRAM_BOT_URL"))
		log.Fatal(err)
	}

	bot.Debug = true

	updateConfig := tgbot.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		cmd := update.Message.Command()
		msg := update.Message

		switch cmd {
		case "q":
			quote := fmt.Sprintf("User: %s\nMessage: %s", msg.From.FirstName, msg.ReplyToMessage.Text)

			msg := tgbot.NewMessage(update.Message.Chat.ID, quote)

			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		}
	}
}

package main

import (
	"github.com/hbourgeot/quoteme/tgbot"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
)

func main() {
	bot, err := tgbot.NewBotAPI(os.Getenv("TELEGRAM_BOT_URL"))
	if err != nil {
		log.Fatal("linea 51", err)
	}

	updateConfig := tgbot.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() || update.Message.ReplyToMessage == nil {
			continue
		}

		cmd := update.Message.Command()

		switch cmd {
		case "q":
			if err = generateQuote(bot, update); err != nil {
				log.Fatal(err)
			}
		}
	}
}

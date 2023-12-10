package main

import (
	"github.com/hbourgeot/quoteme/generators"
	"github.com/hbourgeot/quoteme/tdlib"
	"github.com/hbourgeot/quoteme/tgbot"
	"log"
	"os"
	"strconv"
	"strings"
)

func generateQuote(bot *tgbot.BotAPI, update tgbot.Update) error {
	keyboard := tgbot.NewInlineKeyboardMarkup(
		tgbot.NewInlineKeyboardRow(
			tgbot.NewInlineKeyboardButtonData("👍", "vote_up"),
			tgbot.NewInlineKeyboardButtonData("👎", "vote_down"),
		),
	)

	msg := update.Message

	var count int
	//var messages []string

	initialMessage := msg.ReplyToMessage.MessageID

	if strings.Contains(msg.Text, " ") {
		count, _ = strconv.Atoi(strings.Split(msg.Text, " ")[1])
	} else {
		count = 1
	}

	if count > 1 {
		var messagesIDs []int
		for i := 1; i <= count; i++ {
			messagesIDs = append(messagesIDs, initialMessage)
			initialMessage++
		}

		if err := tdlib.GetMessages(messagesIDs); err != nil {
			return err
		}
	}

	personQuote := msg.ReplyToMessage.From.FirstName + " " + msg.ReplyToMessage.From.LastName
	personQuoteID := msg.ReplyToMessage.From.ID
	photosConf := tgbot.NewUserProfilePhotos(personQuoteID)
	photos, err := bot.GetUserProfilePhotos(photosConf)
	if err != nil {
		return err
	}

	var profile tgbot.PhotoSize

	if photos.TotalCount > 0 {
		profile = photos.Photos[0][0]
	}

	fileConf := tgbot.FileConfig{FileID: profile.FileID}

	profileImg := getImage(bot, fileConf)

	imgBytes, err := generators.GenerateImage(profileImg, personQuote, msg.ReplyToMessage.Text)
	if err != nil {
		log.Fatal("linea 89", err)
	}

	err = os.WriteFile("quote.webp", imgBytes, 0644)
	if err != nil {
		log.Fatal("linea 94", err)
	}

	stkConfig := tgbot.NewSticker(msg.Chat.ID, tgbot.FilePath("quote.webp"))
	stkConfig.ReplyMarkup = keyboard
	stkConfig.ReplyToMessageID = msg.MessageID

	if _, err := bot.Send(stkConfig); err != nil {
		return err
	}

	return nil
}

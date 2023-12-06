package main

import (
	"bytes"
	"fmt"
	"github.com/hbourgeot/quoteme/generators"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/hbourgeot/quoteme/tgbot"
	_ "github.com/joho/godotenv/autoload"
)

func getImageFormat(r io.Reader) (string, error) {
	buf := make([]byte, 512) // 512 bytes should be enough for the magic number
	if _, err := r.Read(buf); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buf)
	switch contentType {
	case "image/jpeg":
		return "jpeg", nil
	case "image/png":
		return "png", nil
	// add other cases as needed
	default:
		return "", fmt.Errorf("unrecognized image format")
	}
}

func decodeImage(r io.Reader, format string) (image.Image, error) {
	switch format {
	case "jpeg":
		return jpeg.Decode(r)
	case "png":
		return png.Decode(r)
	// add other cases as needed
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func getImage(bot *tgbot.BotAPI, fileConf tgbot.FileConfig) image.Image {
	file, err := bot.GetFile(fileConf)
	if err != nil {
		log.Fatal("linea 18", err)
	}

	url, err := bot.GetFileDirectURL(file.FileID)
	if err != nil {
		log.Fatal("linea 23", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("linea 28", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("bad status: %s", resp.Status)
	}

	// Read the entire response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error leyendo")
		return nil
	}
	defer resp.Body.Close()

	// Determine the image format
	format, err := getImageFormat(bytes.NewReader(data))
	if err != nil {
		log.Println("Error obteniendo el formato")
		return nil
	}

	// Decode the image
	img, err := decodeImage(bytes.NewReader(data), format)
	if err != nil {
		log.Println("Error decodificando")
		return nil
	}

	return img
}

func main() {
	bot, err := tgbot.NewBotAPI(os.Getenv("TELEGRAM_BOT_URL"))
	if err != nil {
		log.Fatal("linea 51", err)
	}

	bot.Debug = true

	updateConfig := tgbot.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	keyboard := tgbot.NewInlineKeyboardMarkup(
		tgbot.NewInlineKeyboardRow(
			tgbot.NewInlineKeyboardButtonData("👍", "vote_up"),
			tgbot.NewInlineKeyboardButtonData("👎", "vote_down"),
		),
	)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() && update.Message.ReplyToMessage == nil {
			continue
		}

		cmd := update.Message.Command()

		switch cmd {
		case "q":

			msg := update.Message

			if update.Message.ReplyToMessage == nil {
				continue
			}

			personQuote := msg.ReplyToMessage.From.FirstName + " " + msg.ReplyToMessage.From.LastName
			personQuoteID := msg.ReplyToMessage.From.ID
			photosConf := tgbot.NewUserProfilePhotos(personQuoteID)
			photos, err := bot.GetUserProfilePhotos(photosConf)
			if err != nil {
				log.Fatal("linea 76", err)
			}

			var profile tgbot.PhotoSize

			if photos.TotalCount > 0 {
				profile = photos.Photos[0][0]
			}

			fileConf := tgbot.FileConfig{FileID: profile.FileID}

			profileImg := getImage(bot, fileConf)
			if profileImg == nil {
				continue
			}

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
				log.Fatal("linea 100", err)
			}
		}
	}
}

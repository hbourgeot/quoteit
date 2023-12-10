package tdlib

import (
	"fmt"
	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/gotd/td/tg"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

func initClient() (*gotgproto.Client, error) {
	clientType := gotgproto.ClientType{
		BotToken: os.Getenv("TELEGRAM_BOT_URL"),
	}

	appId, _ := strconv.Atoi(os.Getenv("APP_ID"))
	apiHash := os.Getenv("API_HASH")

	client, err := gotgproto.NewClient(appId, apiHash, clientType, &gotgproto.ClientOpts{Session: sessionMaker.NewInMemorySession(":memory:", 1)})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetMessages(messagesIDs []int) error {
	client, err := initClient()
	if err != nil {
		return err
	}

	ctx := client.CreateContext()
	tgClient := client.API()

	inputMessages := make([]tg.InputMessageClass, 0, len(messagesIDs))
	for _, id := range messagesIDs {
		inputMessageID := &tg.InputMessageID{ID: id}
		inputMessages = append(inputMessages, inputMessageID)
	}

	msgs, err := tgClient.MessagesGetMessages(ctx.Context, inputMessages)
	if err != nil {
		return err
	}

	var messages []tg.MessageClass

	switch m := msgs.(type) {
	case *tg.MessagesMessages:
		messages = m.GetMessages()
		for _, message := range messages {
			switch msg := message.(type) {
			case *tg.Message:
				fmt.Println(msg.Message)

				from := msg.FromID
				switch u := from.(type) {
				case *tg.PeerUser:
					fmt.Println(u.String())
					//user := tg.InputUserClass()
					//tgClient.UsersGetFullUser(ctx.Context, )
				}
			}
		}
	}

	return nil
}

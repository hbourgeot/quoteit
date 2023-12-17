package tdlib

import (
	"context"
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
	var users []tg.PeerClass

	switch m := msgs.(type) {
	case *tg.MessagesMessages:

		messages = m.GetMessages()
		for _, message := range messages {
			switch msg := message.(type) {
			case *tg.Message:

				users = append(users, msg.FromID)
			default:
				continue
			}
		}
	default:
		fmt.Println(m.TypeName())
	}

	fmt.Println("67")
	usersPerMessage, err := getUser(users, tgClient, ctx.Context)
	if err != nil {
		return err
	}

	fmt.Println(len(usersPerMessage))

	return nil
}

func getUser(users []tg.PeerClass, client *tg.Client, ctx context.Context) ([]*tg.User, error) {
	fmt.Println("dentro")
	var usersArray []*tg.User
	var inputUsers []tg.InputUserClass

	for _, user := range users {
		fmt.Println(user.TypeName())
		switch u := user.(type) {
		case *tg.PeerUser:
			inputUsers = append(inputUsers, &tg.InputUser{UserID: u.GetUserID()})
		default:
			continue
		}
	}

	usersClass, err := client.UsersGetUsers(ctx, inputUsers)
	if err != nil {
		return nil, err
	}

	fmt.Println(len(usersClass))

	for _, user := range usersClass {
		switch u := user.(type) {
		case *tg.User:
			usersArray = append(usersArray, u)
		case *tg.UserEmpty:
			fmt.Println("sam")
		default:
			fmt.Println(u.TypeName(), "s")
			continue
		}
	}

	return usersArray, nil
}

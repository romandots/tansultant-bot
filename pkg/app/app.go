package app

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"net/http"
	"strconv"
	"tansulbot/pkg/telegram"
	"time"
)

type App struct {
	db             *badger.DB
	telegramClient *telegram.Client
	Errors         chan error
}

func InitApp(telegramToken string) (*App, error) {
	var err error

	a := &App{}
	a.db, err = badger.Open(badger.DefaultOptions("/tmp/badger"))
	a.telegramClient = initClient(telegramToken, a)

	return a, err
}

func initClient(token string, a telegram.AppInterface) *telegram.Client {
	fmt.Println("Creating client")
	client := &telegram.Client{
		App:           a,
		Token:         token,
		HttpClient:    &http.Client{},
		Commands:      make(map[string]telegram.Command),
		Conversations: make(map[int]*telegram.Conversation),
		Errors:        make(chan error, 100),
	}

	client.Commands["/start"] = telegram.Command{Command: "/start", Description: "Приветствие", Handle: client.CommandWelcome}
	client.Commands["/stop"] = telegram.Command{Command: "/stop", Description: "Конец игры!", Handle: client.CommandGoodbye}

	fmt.Println(client)

	return client
}

func (a *App) Close() error {
	return a.db.Close()
}

func (a *App) ReadTelegramUpdates() {

	fmt.Println("Waiting for updates...")

	go func() {
		fmt.Println("Waiting for telegram errors...")
		for telegramError := range a.telegramClient.Errors {
			if telegramError != nil {
				fmt.Println("Telegram client error:", telegramError.Error())
				a.Errors <- telegramError
			}
		}
	}()

	for {
		updates, err := a.telegramClient.GetUpdates()
		if err != nil {
			fmt.Println("Error getting updates:", err)
		}

		a.telegramClient.ReadUpdates(updates.Result)

		time.Sleep(1 * time.Second)
	}
}

func (a *App) WaitForErrors() {
	fmt.Println("Waiting for app errors...")
	for err := range a.Errors {
		if err != nil {
			fmt.Println("App error:", err.Error())
		}
	}
}

func (a *App) SaveConversation(conversation *telegram.Conversation) error {
	return a.db.Update(func(txn *badger.Txn) error {
		value, err := json.Marshal(*conversation)
		if nil != err {
			return err
		}

		return txn.Set([]byte(strconv.Itoa(conversation.ChatId)), value)
	})
}

// GetConversation returns conversation, existence flag and error
func (a *App) GetConversation(chatId int) (*telegram.Conversation, bool, error) {
	conv := &telegram.Conversation{}
	err := a.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(strconv.Itoa(chatId)))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				// Conversation not found
				return nil
			}
			return err
		}

		var value []byte
		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(value, conv)
		if err != nil {
			return err
		}

		return nil
	})

	if err == nil && conv.ChatId == 0 {
		// not found
		return nil, false, err
	}

	return conv, true, err
}

func (a *App) DropAll() error {
	return a.db.DropAll()
}

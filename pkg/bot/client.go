package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	telegramBaseURL = "https://api.telegram.org/bot"
)

func NewClient(token string) *Client {
	fmt.Println("Creating client")
	client := &Client{
		token:         token,
		httpClient:    &http.Client{},
		commands:      make(map[string]Command),
		conversations: make(map[int]*Conversation),
		Errors:        make(chan error, 100),
	}

	client.commands["/start"] = Command{"/start", "Приветствие", client.CommandWelcome}
	client.commands["/stop"] = Command{"/stop", "Конец игры!", client.CommandGoodbye}

	fmt.Println(client)

	return client
}

func (c *Client) getUrl(endpoint string) string {
	return telegramBaseURL + c.token + endpoint
}

func (c *Client) request(endpoint string) ([]byte, error) {
	url := c.getUrl(endpoint)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) requestf(format string, a ...interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf(format, a...)
	//fmt.Println(endpoint)
	return c.request(endpoint)
}

func (c *Client) interpretCommand(message string) (*Command, bool) {
	command, exists := c.commands[message]

	return &command, exists
}

func (c *Client) runCommand(command *Command, message *Message) (string, error) {
	fmt.Printf("Выполняем команду: %s\n", command.Command)
	reply, err := command.Handle(message)

	if err != nil {
		c.Errors <- err
	}

	return reply, err
}

func (c *Client) sendMessage(chat *Chat, reply *Reply) error {
	var replyMarkupJson []byte
	var err error
	if reply.ReplyMarkup.Keyboard != nil {
		replyMarkupJson, err = json.Marshal(reply.ReplyMarkup)
		if err != nil {
			return err
		}
	}

	fmt.Println("Отправляем сообщение в чат", reply, chat)
	_, err = c.requestf(
		"/sendMessage?chat_id=%d&text=%s&reply_markup=%s",
		chat.Id,
		url.QueryEscape(reply.Message),
		url.QueryEscape(string(replyMarkupJson)),
	)

	c.Errors <- err
	return err
}

func (c *Client) GetUpdates() (*TelegramResponse, error) {
	response, err := c.requestf("/getUpdates?offset=%d", c.lastUpdateId+1)
	if err != nil {
		return nil, err
	}
	var tr TelegramResponse
	err = json.Unmarshal(response, &tr)

	return &tr, err
}

func (c *Client) ReadUpdates(updates []Update) {
	for _, update := range updates {
		var reply string
		var error error
		c.lastUpdateId = update.UpdateId

		fmt.Println("Получено сообщение", update)

		// Act according to conversation state
		conversationState := c.getConversationState(&update.Message.Chat)
		switch conversationState {
		case Unauthorized:
			c.CommandWelcome(&update.Message)
		case Authorization:
			reply, error = c.CommandAuthorize(&update.Message)
			if error != nil {
				c.CommandWelcome(&update.Message)
			}
		default:
			command, exists := c.interpretCommand(update.Message.Text)
			fmt.Println(command, exists)
			if exists {
				fmt.Printf("Обнаружена команда: %s\n", command.Description)
				reply, error = c.runCommand(command, &update.Message)
			}
		}

		// Reply to the message
		if reply != "" {
			c.SendMessage(&update.Message.Chat, reply)
		}
	}
}

func (c *Client) SendMessage(chat *Chat, message string) {
	go c.sendMessage(chat, &Reply{Message: message})
}

func (c *Client) RequestPhoneNumber(chat *Chat) {
	replyMarkup := ReplyMarkup{
		Keyboard: [][]KeyboardButton{
			{{Text: "Отправить мой номер телефона", RequestContact: true}},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	message := "Для авторизации, необходимо, чтобы номер телефона, к которому привязан ваш аккаунт в Телеграме, совпадал " +
		"с номером, который вы указали при регистрации в ШТБП."

	go c.sendMessage(chat, &Reply{Message: message, ReplyMarkup: replyMarkup})
}

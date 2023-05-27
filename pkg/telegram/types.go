package telegram

import (
	"net/http"
)

type AppInterface interface {
	SaveConversation(conversation *Conversation) error
	GetConversation(userId int) (*Conversation, bool, error)
}

type ConversationState int

type Conversation struct {
	ChatId      int
	UserId      int
	PhoneNumber string
	State       ConversationState
}

type Client struct {
	App           AppInterface
	Token         string
	HttpClient    *http.Client
	Commands      map[string]Command
	Conversations map[int]*Conversation
	Errors        chan error
	lastUpdateId  int
}

type Command struct {
	Command     string
	Description string
	Handle      func(message *Message) (string, error)
}

type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
	RequestPoll     *struct {
		Type string `json:"type,omitempty"`
	} `json:"request_poll,omitempty"`
}

type Reply struct {
	Message     string
	ReplyMarkup ReplyMarkup
}

type ReplyMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
}

type Chat struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserId      int    `json:"user_id"`
	Vcard       string `json:"vcard"`
}

type Message struct {
	MessageId int `json:"message_id"`
	From      struct {
		Id           int    `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
	} `json:"from"`
	Chat    Chat    `json:"chat"`
	Date    int     `json:"date"`
	Text    string  `json:"text"`
	Contact Contact `json:"contact"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type TelegramResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

package bot

import (
	"errors"
	"fmt"
)

type ConversationState int

const (
	Unauthorized ConversationState = iota
	Idle
	Authorization
)

type Conversation struct {
	ChatID      int
	PhoneNumber string
	State       ConversationState
}

func (c *Client) createConversation(chat *Chat) *Conversation {
	fmt.Println("Создаем диалог...")
	return &Conversation{
		ChatID: chat.Id,
		State:  Unauthorized,
	}
}

func (c *Client) getConversation(chat *Chat) *Conversation {
	conversation, ok := c.conversations[chat.Id]
	if !ok {
		conversation = c.createConversation(chat)
		c.conversations[chat.Id] = conversation
	}
	fmt.Println(conversation)
	return conversation
}

func (c *Client) setConversationState(chat *Chat, state ConversationState) {
	fmt.Println("Меняем статус диалога...", state)
	conversation := c.getConversation(chat)
	conversation.State = state
	fmt.Println(conversation)
}

func (c *Client) getConversationState(chat *Chat) ConversationState {
	conversation := c.getConversation(chat)
	return conversation.State
}

func (c *Client) attachUserPhoneNumberToConversation(chat *Chat, contact *Contact) error {
	if contact.PhoneNumber == "" {
		return errors.New("Не обнаружили номер телефона")
	}
	fmt.Println("Сохраняем номер телефона в диалог...")
	conversation := c.getConversation(chat)
	conversation.PhoneNumber = contact.PhoneNumber
	fmt.Println(conversation)

	return nil
}

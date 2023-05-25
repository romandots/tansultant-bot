package telegram

import (
	"errors"
	"fmt"
)

const (
	Unauthorized ConversationState = iota
	Idle
	Authorization
)

func (c *Client) createConversation(message *Message) *Conversation {
	fmt.Println("Создаем диалог...")
	return &Conversation{
		ChatId: message.Chat.Id,
		UserId: message.Contact.UserId,
		State:  Unauthorized,
	}
}

func (c *Client) getConversation(message *Message) *Conversation {
	conversation, ok := c.Conversations[message.Contact.UserId]
	if !ok {
		conversation = c.createConversation(message)
		c.Conversations[message.Contact.UserId] = conversation
	}
	fmt.Println(conversation)
	return conversation
}

func (c *Client) setConversationState(message *Message, state ConversationState) {
	fmt.Println("Меняем статус диалога...", state)
	conversation := c.getConversation(message)
	conversation.State = state
	fmt.Println(conversation)
}

func (c *Client) getConversationState(message *Message) ConversationState {
	conversation := c.getConversation(message)
	return conversation.State
}

func (c *Client) attachUserPhoneNumberToConversation(message *Message) error {
	if message.Contact.PhoneNumber == "" {
		return errors.New("Не обнаружили номер телефона")
	}
	fmt.Println("Сохраняем номер телефона в диалог...")
	conversation := c.getConversation(message)
	conversation.PhoneNumber = message.Contact.PhoneNumber
	fmt.Println(conversation)

	return nil
}

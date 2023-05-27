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

// return: conversation, isNew, error
func (c *Client) getConversation(message *Message) (*Conversation, bool, error) {
	fmt.Println("Getting conversation")
	fmt.Printf("Chat ID: %d\n", message.Chat.Id)
	fmt.Printf("%v", message)
	conversation, ok := c.Conversations[message.Chat.Id]

	if !ok {
		fmt.Println("No conversation in the map. Loading from storage...")

		var err error
		conversation, ok, err = c.loadConversation(message.Chat.Id)
		fmt.Printf("%v %v %v", conversation, ok, err)
		if nil != err {
			fmt.Println("Failed loading conversation from storage...")
			return nil, false, err
		}
	}

	if !ok {
		fmt.Println("Didn't find anything in storage. Creating conversation...")
		conversation = c.createConversation(message)
		c.Conversations[message.Chat.Id] = conversation
		return conversation, true, nil
	}

	fmt.Println("Storing conversation in map...")
	c.Conversations[message.Chat.Id] = conversation
	return conversation, false, nil
}

func (c *Client) setConversationState(message *Message, state ConversationState) error {
	fmt.Println("Меняем статус диалога...", state)
	conversation, _, err := c.getConversation(message)
	if err != nil {
		return err
	}
	conversation.State = state
	fmt.Println(conversation)

	if err := c.saveConversation(conversation); err != nil {
		return err
	}

	return nil
}

func (c *Client) getConversationState(message *Message) (ConversationState, error) {
	conversation, isNew, err := c.getConversation(message)
	if err != nil {
		return 0, err
	}

	if isNew {
		return Unauthorized, nil
	}

	return conversation.State, nil
}

func (c *Client) attachUserPhoneNumberToConversation(message *Message) error {
	if message.Contact.PhoneNumber == "" {
		fmt.Print(message)
		return errors.New("Не обнаружили номер телефона")
	}
	fmt.Println("Сохраняем номер телефона в диалог...")
	conversation, _, err := c.getConversation(message)
	if err != nil {
		return err
	}

	conversation.PhoneNumber = message.Contact.PhoneNumber
	fmt.Println(conversation)

	if err := c.saveConversation(conversation); err != nil {
		return err
	}

	return nil
}

func (c *Client) saveConversation(conversation *Conversation) error {
	fmt.Println("Saving conversation")
	return c.App.SaveConversation(conversation)
}

// returns: conversation, exist, error
func (c *Client) loadConversation(userId int) (*Conversation, bool, error) {
	fmt.Println("Loading conversation")
	return c.App.GetConversation(userId)
}

package bot

import "fmt"

func (c *Client) CommandWelcome(message *Message) (string, error) {
	fmt.Println("Выполняем команду /start")

	fmt.Println("Запрашиваем номер телефона...")
	c.RequestPhoneNumber(&message.Chat)

	fmt.Println("Переводим диалог в статус авторизации...")
	c.setConversationState(&message.Chat, Authorization)

	return "", nil
}

func (c *Client) CommandAuthorize(message *Message) (string, error) {
	fmt.Println("Сохраняем номер телефона в диалог...")
	error := c.attachUserPhoneNumberToConversation(&message.Chat, &message.Contact)
	if error != nil {
		return "", error
	}

	c.setConversationState(&message.Chat, Idle)
	conversation := c.getConversation(&message.Chat)

	return "Ура! Теперь мы знаем, кто ты! " + conversation.PhoneNumber, nil
}

func (c *Client) CommandGoodbye(message *Message) (string, error) {
	return "Пока-пока!", nil
}

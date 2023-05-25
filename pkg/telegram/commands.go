package telegram

import "fmt"

func (c *Client) CommandWelcome(message *Message) (string, error) {
	fmt.Println("Выполняем команду /start")

	fmt.Println("Запрашиваем номер телефона...")
	c.RequestPhoneNumber(&message.Chat)

	fmt.Println("Переводим диалог в статус авторизации...")
	c.setConversationState(message, Authorization)

	return "", nil
}

func (c *Client) CommandAuthorize(message *Message) (string, error) {
	fmt.Println("Сохраняем номер телефона в диалог...")
	error := c.attachUserPhoneNumberToConversation(message)
	if error != nil {
		return "", error
	}

	c.setConversationState(message, Idle)
	conversation := c.getConversation(message)

	return "Ура! Теперь мы знаем, кто ты! " + conversation.PhoneNumber, nil
}

func (c *Client) CommandGoodbye(message *Message) (string, error) {
	return "Пока-пока!", nil
}

package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	telegramBaseURL = "https://api.telegram.org/bot"
	token           = "5644445038:AAFTsJWu8JSITroIN_poVLIljojgaDMVXmc"
)

type Client struct {
	token      string
	httpClient *http.Client
}

func (c Client) GetUrl(endpoint string) string {
	return telegramBaseURL + c.token + endpoint
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
	Chat struct {
		Id        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		Type      string `json:"type"`
	} `json:"chat"`
	Date int    `json:"date"`
	Text string `json:"text"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type TelegramResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{},
	}
}

func HandleError(e error) {
	log.Fatal(e)
}

func (c Client) Request(endpoint string) ([]byte, error) {
	url := c.GetUrl(endpoint)
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

func (c Client) GetUpdates() (*TelegramResponse, error) {
	response, err := c.Request("/getUpdates")
	if err != nil {
		return nil, err
	}
	var tr TelegramResponse
	json.Unmarshal(response, &tr)

	return &tr, nil
}

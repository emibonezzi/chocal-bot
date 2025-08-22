package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Bot struct {
	telegramEndpoint string
	botToken         string
	currentUser      User
}

type User struct {
	id        int
	firstName string
	lastName  string
	username  string
}

type NewMessageUpdate struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

// populate bot with useful and reusable values related to the user that is interacting with the bot
func InitializeBot(body string) (Bot, error) {
	var b Bot
	b.telegramEndpoint = os.Getenv("API_ENDPOINT")
	b.botToken = os.Getenv("BOT_TOKEN")
	var m NewMessageUpdate
	err := json.Unmarshal([]byte(body), &m)
	if err != nil {
		log.Printf("Failed to parse telegram update: %v", err)
		return b, err
	}

	b.currentUser = User{
		id:        m.Message.From.ID,
		firstName: m.Message.From.FirstName,
		lastName:  m.Message.From.LastName,
		username:  m.Message.From.Username,
	}

	return b, err

}

func (b *Bot) SendText(text string, chatID int) (*http.Response, error) {
	type RequestBody struct {
		ChatID int    `json:"chat_id"`
		Text   string `json:"text"`
	}
	bs := RequestBody{
		ChatID: chatID,
		Text:   text,
	}

	body, err := json.Marshal(&bs)
	if err != nil {
		log.Printf("Error in marshaling body: %v", err)
		return &http.Response{}, err
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", b.telegramEndpoint, b.botToken)
	payload := bytes.NewReader(body)
	return http.Post(url, "application/json", payload)
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Bot struct {
	Endpoint string
	Token    string
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

	url := fmt.Sprintf("%s/bot%s/sendMessage", b.Endpoint, b.Token)
	payload := bytes.NewReader(body)
	return http.Post(url, "application/json", payload)
}

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

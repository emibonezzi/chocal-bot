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
	currentMessage   CurrentMessage
}

type User struct {
	id        int
	firstName string
	lastName  string
	username  string
}

type CurrentMessage struct {
	text string
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

	b.currentMessage = CurrentMessage{
		text: m.Message.Text,
	}

	return b, err

}

func (b *Bot) GreetUser() (*http.Response, error) {
	type RequestBody struct {
		ChatID int    `json:"chat_id"`
		Text   string `json:"text"`
	}

	greetingText := fmt.Sprintf("Hello %s! Welcome to ChoCal bot. ChoCal is a customizable Telegram Bot to get daily notifications about your favorite football ⚽ team(s)' fixtures. Use /list to see the list of available teams.", b.currentUser.firstName)

	bs := RequestBody{
		ChatID: b.currentUser.id,
		Text:   greetingText,
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

func (b *Bot) DisplayTeams(list TeamsList) (*http.Response, error) {
	type InlineKeyboardButton struct {
		Text         string `json:"text"`
		CallbackData string `json:"callback_data"`
	}
	type InlineKeyboardMarkup struct {
		InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
	}

	type RequestBody struct {
		ChatID      int                  `json:"chat_id"`
		Text        string               `json:"text"`
		ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
	}

	buttons := make([]InlineKeyboardButton, 0, len(list.Teams))
	mainArray := make([][]InlineKeyboardButton, 0, 1)

	for _, t := range list.Teams {
		b := InlineKeyboardButton{
			Text:         fmt.Sprintf("%s %s", t.Emoji, t.Name),
			CallbackData: t.Id,
		}
		buttons = append(buttons, b)
	}

	mainArray = append(mainArray, buttons)
	inlineKeyboard := InlineKeyboardMarkup{
		InlineKeyboard: mainArray,
	}

	body := RequestBody{
		ChatID:      b.currentUser.id,
		Text:        "Please select your favorite team to start receving notifications.",
		ReplyMarkup: inlineKeyboard,
	}

	bs, err := json.Marshal(&body)
	if err != nil {
		log.Printf("Error in parsing reply markup body: %v", err)
		return &http.Response{}, err
	}

	fmt.Print(string(bs))

	url := fmt.Sprintf("%s/bot%s/sendMessage", b.telegramEndpoint, b.botToken)
	payload := bytes.NewReader(bs)
	return http.Post(url, "application/json", payload)

}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := loadConfig()

	var m NewMessageUpdate
	err := json.Unmarshal([]byte(req.Body), &m)
	if err != nil {
		log.Printf("Failed to unmarshal update: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Fail to unmarshal update",
		}, err
	}

	if m.Message.Text == "/start" {
		// send text
		// create body struct and populate
		type RequestBody struct {
			ChatID int    `json:"chat_id"`
			Text   string `json:"text"`
		}
		b := RequestBody{
			ChatID: m.Message.Chat.ID,
			Text:   fmt.Sprintf("Hello @%v", m.Message.From.Username),
		}

		// convert it to json
		rb, err := json.Marshal(&b)
		if err != nil {
			log.Printf("Failed to marshal request body: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Fail to marshal request body",
			}, err
		}

		// create http request
		response, err := http.Post(c.apiEndpoint, "application/json", bytes.NewReader(rb))
		if err != nil {
			log.Printf("Failed to create http request: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Fail to make request to telegram servers",
			}, err
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			log.Printf("Fail to make request to telegram servers: %v", response.Status)
			return events.APIGatewayProxyResponse{
				StatusCode: response.StatusCode,
				Body:       "Fail to make request to telegram servers",
			}, err

		}

	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Success",
	}, nil
}

func main() {
	lambda.Start(handler)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// read teams.json from S3
	c, err := loadConfig(ctx) // pass context object (that includes AWS creds and more) so config can read teams from S3
	if err != nil {
		log.Printf("Failed to read teams.json: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to read teams.json",
		}, err
	}

	fmt.Print(c.teamsList)

	b, err := InitializeBot(req.Body)
	if err != nil {
		log.Printf("Failed to initialize bot: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to initialize Telegram Bot"}, err
	}

	// parse update coming from Telegram
	var m NewMessageUpdate
	err = json.Unmarshal([]byte(req.Body), &m)
	if err != nil {
		log.Printf("Failed to unmarshal update: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to unmarshal update",
		}, err
	}

	// greet user
	if m.Message.Text == "/start" {
		message, err := b.SendText(fmt.Sprintf("Hello @%v", m.Message.From.Username), m.Message.Chat.ID)
		if err != nil {
			log.Printf("Failed to send text: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Failed to send text",
			}, err
		}

		if message.StatusCode != 200 {
			log.Printf("Telegram returned an error: %v", message.StatusCode)
			return events.APIGatewayProxyResponse{
				StatusCode: message.StatusCode,
				Body:       "Telegram return an error",
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

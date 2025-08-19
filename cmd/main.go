package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load env vars: %v", err)
	}

	c := loadConfig()

	var m NewMessageUpdate
	err = json.Unmarshal([]byte(req.Body), &m)
	if err != nil {
		log.Printf("Failed to unmarshal update: %v", err)
		return events.APIGatewayProxyResponse{}, err
	}

	if m.Message.Text == "/start" {
		// send text
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello from Go Lambda!",
	}, nil
}

func sendText(text string) {

}

func main() {
	lambda.Start(handler)
}

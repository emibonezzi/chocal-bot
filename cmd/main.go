package main

import (
	"context"
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
			StatusCode: 200, // to avoid telegram server retry
			Body:       "Failed to read teams.json",
		}, err
	}

	// load bot
	b, err := InitializeBot(req.Body)
	if err != nil {
		log.Printf("Failed to initialize bot: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Failed to initialize Telegram Bot"}, err
	}

	// init db
	db, err := LoadClient(ctx)
	defer db.DisconnectClient(ctx)

	// greet user
	if b.currentMessage.text == "/start" {
		message, err := b.GreetUser()
		if err != nil {
			log.Printf("Failed to send text: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
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

		// save user in db
		_, err = db.SaveUser(ctx, b.currentUser.id)
		if err != nil {
			log.Printf("Error in saving user in db: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       "Error in saving user in DB"}, err
		}
	}

	// display teams to user
	if b.currentMessage.text == "/list" {
		message, err := b.DisplayTeams(c.teamsList)
		if err != nil {
			log.Printf("Error in displaying teams: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       "Error in displaying teams",
			}, err
		}

		if message.StatusCode != 200 {
			log.Printf("Telegram returned an error: %v", message.Status)
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       "Telegram returned an error",
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

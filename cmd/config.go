package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	apiEndpoint string
	botToken    string
	teamsList   []byte
}

type TeamsList struct {
	Teams []Team `json:"teams"`
}

type Team struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
}

func loadConfig(ctx context.Context) (Config, error) {
	var c Config
	// load env vars
	c.apiEndpoint = os.Getenv("API_ENDPOINT")
	c.botToken = os.Getenv("BOT_TOKEN")

	// load teams list from s3
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Error in loading lambda context: %v", err)
		return Config{}, err
	}
	client := s3.NewFromConfig(cfg)
	resp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("all-teams-bucket"),
		Key:    aws.String("teams.json"),
	})
	if err != nil {
		log.Printf("Error in reading teams.json from S3: %v", err)
		return Config{}, err
	}

	defer resp.Body.Close()

	teams, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error in parsing teams.json: %v", err)
		return Config{}, err
	}

	c.teamsList = teams

	return c, err
}

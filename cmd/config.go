package main

import "os"

type Config struct {
	apiEndpoint string
	botToken    string
}

func loadConfig() Config {
	var c Config
	c.apiEndpoint = os.Getenv("API_ENDPOINT")
	c.botToken = os.Getenv("BOT_TOKEN")
	return c
}

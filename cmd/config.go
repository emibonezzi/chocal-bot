package main

import "os"

type Config struct {
	botToken string
}

func loadConfig() Config {
	var c Config
	c.botToken = os.Getenv("BOT_TOKEN")
	return c
}

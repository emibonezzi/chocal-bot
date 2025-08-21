package main

import "os"

type Config struct {
	apiEndpoint string
}

func loadConfig() Config {
	var c Config
	c.apiEndpoint = os.Getenv("API_ENDPOINT")
	return c
}

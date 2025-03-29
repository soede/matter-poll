package config

import (
	"log"
	"os"
)

type Config struct {
	Env               string
	MattermostURL     string
	BotToken          string
	TarantoolHost     string
	TarantoolUser     string
	TarantoolPassword string
}

func MustLoad() *Config {

	env := os.Getenv("APP_ENV")

	mattermostURL := os.Getenv("MATTERMOST_URL")
	if mattermostURL == "" {
		log.Fatal("Failed to find the MATTERMOST_URL environment variable.")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Failed to find the BOT_TOKEN environment variable.")
	}

	tarantoolHost := os.Getenv("TARANTOOL_HOST")
	if tarantoolHost == "" {
		log.Fatal("Failed to find the TARANTOOL_HOST environment variable.")
	}

	tarantoolUser := os.Getenv("TARANTOOL_USER")
	if tarantoolUser == "" {
		log.Fatal("Failed to find the TARANTOOL_USER environment variable.")
	}

	tarantoolPassword := os.Getenv("TARANTOOL_PASSWORD")

	return &Config{
		env,
		mattermostURL,
		botToken,
		tarantoolHost,
		tarantoolUser,
		tarantoolPassword}
}

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type TelegramConfig struct {
	BotToken string
	BotUrl string
}

func LoadTelegramConfig() (*TelegramConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
		return nil, err
	}

	config := &TelegramConfig{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		BotUrl: os.Getenv("TELEGRAM_BOT_URL"),
	}

	return config, nil
}

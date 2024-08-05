package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

const confPath = "internal/config/.env"

type Config struct {
	TelegramBotToken string
	TelegramBotName  string
}

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("environment variable %s not set", key)
}
func GetConfig() (*Config, error) {
	config := &Config{}
	err := godotenv.Load(confPath)
	config.TelegramBotToken, err = getEnv("TELEGRAM_BOT_TOKEN")
	if err != nil {
		return nil, fmt.Errorf("could not parse config from .env:%s", err)
	}
	config.TelegramBotName, err = getEnv("TELEGRAM_BOT_NAME")
	if err != nil {
		return nil, fmt.Errorf("could not parse config from .env:%s", err)
	}
	return config, nil
}

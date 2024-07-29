package server

import (
	tgClient "authservice/internal/clients/telegram"
	"authservice/internal/config"
	"log"
)

var TgClient *tgClient.Client

func InitTelegramClient() {
	cfg := config.GetConfig()

	TgClient = tgClient.New(cfg.TgBotHost, cfg.TelegramBotToken)
	log.Println("Telegram client initialized")
}

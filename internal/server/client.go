package server

import (
	"authservice/config"
	tgClient "authservice/internal/clients/telegram"
	"log"
)

var TgClient *tgClient.Client

func InitTelegramClient() {
	cfg := config.GetConfig()

	TgClient = tgClient.New(cfg.TgBotHost, cfg.TelegramBotToken)
	log.Println("Telegram client initialized")
}

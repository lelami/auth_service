package server

import (
	tgClient "authservice/internal/clients/telegram"
	"log"
)

var TgClient *tgClient.Client

const (
	tgBotHost        = "api.telegram.org"
	batchSize        = 100
	telegramBotToken = "7391650942:AAENvCBC8wTQXnaAFOBWpEQtm7eTR3mxFWw"
)

func InitTelegramClient() {
	TgClient = tgClient.New(tgBotHost, telegramBotToken)
	log.Println("Telegram client initialized")
}

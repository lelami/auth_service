package app

import (
	"context"
	"log"

	tgClient "authservice/internal/clients/telegram"
	event_consumer "authservice/internal/consumer/event-consumer"
	"authservice/internal/events/telegram"
	"authservice/internal/repository/userdb"
)

const (
	tgBotHost        = "api.telegram.org"
	batchSize        = 100
	telegramBotToken = "7391650942:AAENvCBC8wTQXnaAFOBWpEQtm7eTR3mxFWw"
)

func RunTG(ctx context.Context, userDB userdb.DB) {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, telegramBotToken),
		userDB,
	)

	log.Print("INFO tg bot started")

	consumer := event_consumer.New(ctx, eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("INFO tg bot is stopped", err)
	}
}

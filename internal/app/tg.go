package app

import (
	"context"
	"log"

	event_consumer "authservice/internal/consumer/event-consumer"
	"authservice/internal/events/telegram"
	"authservice/internal/repository/userdb"
	"authservice/internal/server"
)

const (
	batchSize = 100
)

func RunTG(ctx context.Context, userDB userdb.DB) {
	server.InitTelegramClient()

	eventsProcessor := telegram.New(
		server.TgClient,
		userDB,
	)

	log.Print("INFO tg bot started")

	consumer := event_consumer.New(ctx, eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("INFO tg bot is stopped", err)
	}
}

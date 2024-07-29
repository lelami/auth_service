package app

import (
	"context"
	"log"

	"authservice/internal/config"
	event_consumer "authservice/internal/consumer/event-consumer"
	"authservice/internal/events/telegram"
	"authservice/internal/repository/userdb"
	"authservice/internal/server"
)

func RunTG(ctx context.Context, userDB userdb.DB) {
	cfg := config.GetConfig()

	server.InitTelegramClient()

	eventsProcessor := telegram.New(
		server.TgClient,
		userDB,
	)

	log.Print("INFO tg bot started")

	consumer := event_consumer.New(ctx, eventsProcessor, eventsProcessor, cfg.BatchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("INFO tg bot is stopped", err)
	}
}

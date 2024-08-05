package app

import (
	"authservice/internal/clients/telegram"
	"authservice/internal/config"
	"authservice/internal/consumer"
	"authservice/internal/handler/httphandler"
	"authservice/internal/repository/cache"
	"authservice/internal/server"
	"authservice/internal/service"
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
)

const (
	tgBotHost = "api.telegram.org"
)

func Run() {

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	// initialize dbs
	userDB, err := cache.UserCacheInit(ctx, &wg)
	if err != nil {
		log.Fatalf("ERROR failed to initialize user database: %v", err)
	}
	tokenDB, err := cache.TokenCacheInit(ctx, &wg)
	if err != nil {
		log.Fatalf("ERROR failed to initialize tokens database: %v", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("ERROR failed to load config file: %s", err)
	}

	tgClient := telegram.New(tgBotHost, cfg.TelegramBotToken)
	// initialize service
	service.Init(userDB, tokenDB, cfg.TelegramBotName)
	consumer.Init(tgClient)
	consumer.StartTelegramConsumer()
	go func() {
		err := server.Run("", "8000", httphandler.NewRouter())
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("ERROR server run ", err)
		}
	}()

	log.Println("INFO auth service is running")

	<-ctx.Done()

	err = server.Shutdown()
	if err != nil {
		log.Fatal("ERROR server was not gracefully shutdown", err)
	}
	wg.Wait()

	log.Println("INFO auth service was gracefully shutdown")
}

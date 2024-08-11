package app

import (
	"authservice/internal/config"
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
	telegramTokenDB, err := cache.TelegramAuthCodeCacheInit(ctx, &wg)
	if err != nil {
		log.Fatalf("ERROR failed to initialize tokens database: %v", err)
	}

	telegramCfg, err := config.LoadTelegramConfig()
	if err != nil {
		log.Fatalf("ERROR failed to load config: %v", err)
	}
	telegramService, err := service.NewTelegramService(telegramCfg)
	if err != nil {
		log.Fatalf("ERROR failed to initialize Telegram service: %v", err)
	}

	// initialize service
	service.Init(userDB, tokenDB, telegramTokenDB, *telegramService)

	go func() {
		err := server.Run("localhost", "8000", httphandler.NewRouter())
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("ERROR server run ", err)
		}
	}()

	go func() {
		err := telegramService.ListenForUpdates(ctx)
		if err != nil {
			log.Fatalf("ERROR telegram bot failed: %v", err)
		}
	}()

	log.Println("INFO auth ice is running")

	<-ctx.Done()

	err = server.Shutdown()
	if err != nil {
		log.Fatal("ERROR server was not gracefully shutdown", err)
	}
	wg.Wait()

	log.Println("INFO auth service was gracefully shutdown")
}

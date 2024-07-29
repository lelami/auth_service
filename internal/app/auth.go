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
	"time"
)

func Run() {

	cfg := config.GetConfig()

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
	otpDB, err := cache.OTPCacheInit(ctx, &wg)
	if err != nil {
		log.Fatalf("ERROR failed to initialize OTPs database: %v", err)
	}

	// initialize service
	service.Init(userDB, tokenDB, otpDB)

	wg.Add(1)
	go func() {
		defer wg.Done()
		otpDB.StartCleanupOTPDaemon(ctx, 1*time.Minute)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.Run(cfg.Host, cfg.Port, httphandler.NewRouter())
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("ERROR server run ", err)
		}
	}()

	// start tg bot
	wg.Add(1)
	go func() {
		defer wg.Done()
		RunTG(ctx, userDB)
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

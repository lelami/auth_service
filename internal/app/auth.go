package app

import (
	"authservice/internal/handler/httphandler"
	"authservice/internal/repository/cache"
	http2 "authservice/internal/server/http"
	"authservice/internal/service"
	"authservice/pkg/meter"
	"authservice/pkg/tracer"
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
)

var serviceName = "Auth Service"

func Run() {

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	traceProv, err := tracer.InitTracer("http://localhost:14268/api/traces", serviceName)
	if err != nil {
		log.Fatal("init tracer", err)
	}

	meterProv, err := meter.InitMeter(ctx, serviceName)
	if err != nil {
		log.Fatal("init meter", err)
	}

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

	// initialize service
	service.Init(userDB, tokenDB)

	go func() {
		err := http2.Run("localhost", "8000", httphandler.NewRouterWithTrace())
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("ERROR server run ", err)
		}
	}()

	log.Println("INFO auth service is running")

	<-ctx.Done()

	err = http2.Shutdown()
	if err != nil {
		log.Fatal("ERROR server was not gracefully shutdown", err)
	}

	if err := traceProv.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	}

	if err := meterProv.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down meter provider: %v", err)
	}

	wg.Wait()

	log.Println("INFO auth service was gracefully shutdown")
}

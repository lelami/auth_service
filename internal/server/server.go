package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	sw "github.com/swaggo/http-swagger"
)

var server *http.Server

func Run(host, port string, router *http.ServeMux) error {
	// comment for debug mode
	server = &http.Server{
		Addr:              fmt.Sprintf("%s:%s", host, port),
		Handler:           http.TimeoutHandler(router, 1*time.Second, ""),
		ReadHeaderTimeout: 200 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
	}

	// uncomment for debug mode
	/*	server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: router,
	}*/

	return server.ListenAndServe()
}

func Shutdown() error {
	return server.Shutdown(context.Background())
}

func ServerDocs(ctx context.Context, host string) error {
	router := mux.NewRouter()
	router.PathPrefix("/docs/").Handler(sw.WrapHandler)

	srv := &http.Server{
		Addr:    host,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("ERROR swagger server run: %v", err)
		}
	}()

	<-ctx.Done()

	return nil
}

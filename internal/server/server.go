package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
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

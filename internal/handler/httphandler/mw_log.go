package httphandler

import (
	"log"
	"net/http"
)

func LogUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		log.Println("Got request", req.Method, req.URL.Path, req.RemoteAddr, req.UserAgent(),
			"user", req.Header.Get(HeaderUserID))
		next.ServeHTTP(resp, req)

	})
}

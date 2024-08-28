package httphandler

import (
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		resp.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, User-ID, X-Service-Key")

		if req.Method == http.MethodOptions {
			resp.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(resp, req)
	})
}

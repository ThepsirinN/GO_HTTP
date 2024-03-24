package middleware

import (
	"go_http_barko/utility/logger"
	"net/http"
)

func MiddlewareOne(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.Context(), "Middleware Execute!")
		next.ServeHTTP(w, r)
	})
}

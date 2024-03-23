package middleware

import (
	"fmt"
	"net/http"
)

func MiddleWareOne(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println()
		fmt.Println("Middleware Execute")
		fmt.Println()
		next.ServeHTTP(w, r)
	})
}

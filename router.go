package main

import (
	"fmt"
	"go_http_barko/middleware"
	"net/http"
)

func initRounter() *http.ServeMux {
	mux := http.NewServeMux()
	healthCheck(mux)
	readinessCheck(mux)
	routerGroupV1(mux)
	return mux
}

func healthCheck(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Write([]byte("Your service is running"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func readinessCheck(mux *http.ServeMux) {
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Write([]byte("Your service is ready"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func routerGroupV1(mux *http.ServeMux) {
	apiV1 := "/api/v1"
	// propagator := propagation.NewCompositeTextMapPropagator(
	// 	propagation.TraceContext{},
	// 	propagation.Baggage{},
	// )

	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Write([]byte("Hello,World!"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}

	mux.Handle(fmt.Sprint(apiV1, "/hello"), middleware.MiddleWareOne(handler))

}

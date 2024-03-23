package main

import (
	"fmt"
	"go_http_barko/middleware"
	"net/http"
)

type handlerV1 interface {
	GetHelloWorldHandler(w http.ResponseWriter, r *http.Request)
}

func initRounter(rv1 handlerV1) *http.ServeMux {
	mux := http.NewServeMux()
	healthCheck(mux)
	readinessCheck(mux)
	routerGroupV1(mux, rv1)
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

func routerGroupV1(mux *http.ServeMux, hv1 handlerV1) {
	apiV1 := "/api/v1"
	// propagator := propagation.NewCompositeTextMapPropagator(
	// 	propagation.TraceContext{},
	// 	propagation.Baggage{},
	// )

	mux.Handle(fmt.Sprint(apiV1, "/hello"), middleware.MiddleWareOne(hv1.GetHelloWorldHandler))

}

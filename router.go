package main

import (
	"fmt"
	"go_http_barko/middleware"
	httpLogger "go_http_barko/utility/httpLogger"
	otelhandle "go_http_barko/utility/otelHandle"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const API_V1 = "/api/v1"

type handlerV1 interface {
	GetAllUserHandler(w http.ResponseWriter, r *http.Request)
}

func initRounter(rv1 handlerV1) http.Handler {
	mux := http.NewServeMux()
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(propagator)
	healthCheck(mux)
	readinessCheck(mux)
	routerGroupV1(mux, rv1)

	return otelhttp.NewHandler(httpLogger.HttpLoggerMiddleware(mux), API_V1)
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
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Your service is ready"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func routerGroupV1(mux *http.ServeMux, hv1 handlerV1) {

	// mux.Handle(fmt.Sprint(apiV1, "/hello"), middleware.MiddleWareOne(hv1.GetHelloWorldHandler))
	otelhandle.OtelHttpHandleFunc(fmt.Sprint(API_V1, "/hello"), mux, middleware.MiddlewareOne(hv1.GetAllUserHandler))

}

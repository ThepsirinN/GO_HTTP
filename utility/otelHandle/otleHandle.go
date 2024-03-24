package otelhandle

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func OtelHttpHandleFunc(pattern string, mux *http.ServeMux, next http.HandlerFunc) {
	handler := otelhttp.WithRouteTag(pattern, next)
	mux.Handle(pattern, handler)
}

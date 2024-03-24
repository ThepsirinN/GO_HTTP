package tracer

import (
	"context"
	"fmt"
	"go_http_barko/utility/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func InitTraceProvider(ctx context.Context, env, appName string) *sdkTrace.TracerProvider {

	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(appName),
	))

	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("resource.Merge %v", err))
	}

	// create trace provider without exporter
	sampler := sdkTrace.ParentBased(
		sdkTrace.AlwaysSample(),
		sdkTrace.WithRemoteParentSampled(sdkTrace.AlwaysSample()),
	)

	var tp *sdkTrace.TracerProvider

	if env == "local" {
		tp = sdkTrace.NewTracerProvider(
			sdkTrace.WithResource(res),
			sdkTrace.WithSampler(sampler),
		)
	}

	// register the global Tracer provider
	otel.SetTracerProvider(tp)
	tracer = otel.GetTracerProvider().Tracer("barko.com/trace")

	return tp
}

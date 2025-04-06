// Package tracer exposes the methods used for otel instrumentation
package tracer

import (
	"context"
	"fmt"
	"os"

	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tp *sdktrace.TracerProvider
var tracer trace.Tracer
var once sync.Once

// InitTracer starts the otel tracer
func InitTracer() (*sdktrace.TracerProvider, error) {
	headers := map[string]string{
		"content-type": "application/json",
	}
	ep := os.Getenv("OTEL_EP")
	fmt.Printf("\tusing OTEL_EP=%s\n", ep)
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(ep),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, err
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "shorturl"),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		return nil, err
	}

	once.Do(func() {
		tp = sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		)
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(propagation.TraceContext{},
				propagation.Baggage{}),
		)

		tracer = tp.Tracer("shorturl")
	})

	return tp, nil
}

func GetTracer() trace.Tracer {
	return tracer
}

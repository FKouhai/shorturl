// Package tracer_test contains tests for the tracer package
package tracer_test

import (
	"context"
	"fmt"
	"os"

	"testing"
	"url_shortener/tracer"

	"github.com/stretchr/testify/assert"
)

func TestInitTracer_Success(t *testing.T) {
	os.Setenv("OTEL_EP", "otelcollector.universe.home:443")
	tp, err := tracer.InitTracer()
	assert.NoError(t, err)
	assert.NotNil(t, tp)
}

func TestInitTracer_NoEndpoint(t *testing.T) {
	os.Unsetenv("OTEL_EP")
	_, err := tracer.InitTracer()
	assert.Error(t, err)
	assert.EqualError(t, err, "Error: OTEL_EP is not set")
}

func TestGetTracer_Success(t *testing.T) {
	os.Setenv("OTEL_EP", "otelcollector.universe.home:443")
	tracer.InitTracer()
	tra := tracer.GetTracer()
	assert.NotNil(t, tra)
}

func TestGetTracer_Error(t *testing.T) {
	os.Unsetenv("OTEL_EP")
	_, err := tracer.InitTracer()
	tr := tracer.GetTracer()
	assert.Error(t, err)
	assert.NotNil(t, tr)
}

func TestTracerUsage(t *testing.T) {
	os.Setenv("OTEL_EP", "otelcollector.universe.home:443")
	tp, _ := tracer.InitTracer()
	tr := tp.Tracer("test-tracer")

	_, span := tr.Start(context.Background(), "TestSpan")
	defer span.End()

	span.AddEvent("TestEvent")
	assert.NotNil(t, span)
}

func TestTracerUsage_NoInit(t *testing.T) {
	os.Unsetenv("OTEL_EP")
	tp, _ := tracer.InitTracer()
	if tp != nil {
		fmt.Println("InitTracer should have failed due to missing OTEL_EP")
		return
	}
	_ = tracer.GetTracer()
	_, span := tracer.GetTracer().Start(context.Background(), "TestSpan")
	defer span.End()

	span.AddEvent("TestEvent")
	assert.NotNil(t, span)
}

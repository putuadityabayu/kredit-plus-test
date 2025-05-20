/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package otel

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"time"
	"xyz/pkg/response"

	"go.opentelemetry.io/otel"
	otels "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otlptrace "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
)

// InitTelemetry initializes OpenTelemetry with trace and metric exporters
func InitTelemetry(ctx context.Context, serviceName string) {
	if env := viper.GetString("app_env"); env == "test" {
		fmt.Println("TEST")
		setupTest(serviceName)
		return
	}

	// Create a resource detailing service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	// ===== TRACES =====
	// Create and configure trace exporter
	traceExporter, err := otlptrace.New(ctx,
		otlptrace.WithEndpoint(viper.GetString("otel_url")), // Pastikan port ini adalah port yang di-publish
		otlptrace.WithInsecure(),
		otlptrace.WithTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create trace exporter: %v", err)
	}

	// Create trace provider with exporter
	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(traceExporter),
	)
	otel.SetTracerProvider(tracerProvider)

	// Get a tracer
	tracer = tracerProvider.Tracer(serviceName)

}

func Shutdown() {
	if tracerProvider != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}
}

// StartSpanHandler is to start parent span, and add http attributes
func StartSpanHandler(c *fiber.Ctx, spanName string) (context.Context, *Span) {
	ctx, span := StartSpan(c.UserContext(), spanName)

	span.SetAttributes(
		attribute.String("http.method", c.Method()),
		attribute.String("http.path", c.Path()),
		attribute.String("ip_address", c.IP()),
	)

	c.SetUserContext(ctx)
	return ctx, &Span{span}
}

// StartSpan starts a new span with the given name and returns the span with context
func StartSpan(ctx context.Context, spanName string) (context.Context, *Span) {
	if tracer == nil {
		return nil, nil
	}
	ctx, span := tracer.Start(ctx, spanName)
	return ctx, &Span{span}
}

// FromContext Get parent span from context
func FromContext(ctx context.Context) *Span {
	span := trace.SpanFromContext(ctx)
	return &Span{span}
}

// Span helper for span
type Span struct {
	trace.Span
}

func ToKeyValue(attributes map[string]any) []attribute.KeyValue {
	var attrs []attribute.KeyValue

	for k, v := range attributes {
		switch val := v.(type) {
		case string:
			attrs = append(attrs, attribute.String(k, val))
		case int:
			attrs = append(attrs, attribute.Int(k, val))
		case int64:
			attrs = append(attrs, attribute.Int64(k, val))
		case float64:
			attrs = append(attrs, attribute.Float64(k, val))
		case bool:
			attrs = append(attrs, attribute.Bool(k, val))
		case []string:
			attrs = append(attrs, attribute.StringSlice(k, val))
		case []int:
			attrs = append(attrs, attribute.IntSlice(k, val))
		case []int64:
			attrs = append(attrs, attribute.Int64Slice(k, val))
		case []float64:
			attrs = append(attrs, attribute.Float64Slice(k, val))
		case []bool:
			attrs = append(attrs, attribute.BoolSlice(k, val))
		}
	}
	return attrs
}

// AddEventHelper helper for add event
func (s *Span) AddEventHelper(name string, attributes map[string]any) {
	if attributes != nil {
		s.AddEvent(name, trace.WithAttributes(ToKeyValue(attributes)...))
	}
}

type ErrorWithStack interface {
	ErrorStack() string
}

// RecordErrorHelper helper for record error
func (s *Span) RecordErrorHelper(err error, message string) {
	var e response.ErrorResponse
	if !errors.As(err, &e) {
		e = response.NewError(0, "", "", nil, err)
	}
	s.RecordError(err, trace.WithAttributes(ToKeyValue(map[string]any{"exception.stacktrace": e.ErrorStack()})...))
	s.SetStatus(codes.Error, message)
}

// GetTraceID to get trace id and span id
func (s *Span) GetTraceID() string {
	sc := s.SpanContext()

	return sc.TraceID().String()
}

func setupTest(serviceName string) {
	// Set global variables
	traceExporter := tracetest.NewInMemoryExporter()
	// Create trace provider with exporter
	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(traceExporter),
	)

	// Step 2: Set global tracer provider
	otels.SetTracerProvider(tracerProvider)

	tracer = tracerProvider.Tracer(serviceName)
}

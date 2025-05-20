/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package otel

import (
	"bytes"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"os"
	"testing"
)

// Mock setup
type MockSpan struct {
	mock.Mock
	trace.Span
}

func (m *MockSpan) End(options ...trace.SpanEndOption) {
	m.Called(options)
}

func (m *MockSpan) AddEvent(name string, options ...trace.EventOption) {
	m.Called(name, options)
}

func (m *MockSpan) RecordError(err error, options ...trace.EventOption) {
	m.Called(err, options)
}

func (m *MockSpan) SetStatus(code codes.Code, description string) {
	m.Called(code, description)
}

func (m *MockSpan) SetAttributes(attributes ...attribute.KeyValue) {
	m.Called(attributes)
}

func (m *MockSpan) SpanContext() trace.SpanContext {
	args := m.Called()
	return args.Get(0).(trace.SpanContext)
}

func (m *MockSpan) IsRecording() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSpan) TracerProvider() trace.TracerProvider {
	args := m.Called()
	return args.Get(0).(trace.TracerProvider)
}

func TestMain(m *testing.M) {
	// Setup
	setupTestEnv()

	// Run tests
	code := m.Run()

	// Teardown
	teardownTestEnv()

	os.Exit(code)
}

func setupTestEnv() {
	config := []byte(`{"app_env":"test"}`)
	viper.SetConfigType("json")
	_ = viper.ReadConfig(bytes.NewBuffer(config))
}

func teardownTestEnv() {
	// Ensure we clean up the tracer provider
	if tracerProvider != nil {
		_ = tracerProvider.Shutdown(context.Background())
	}
}

func TestInitTelemetry(t *testing.T) {
	// Reset global variables
	tracerProvider = nil
	tracer = nil

	// Test initialization
	InitTelemetry(context.Background(), "test-service")

	// Verify tracer provider and tracer are set
	assert.NotNil(t, tracerProvider)
	assert.NotNil(t, tracer)

	// Call Shutdown to clean up
	Shutdown()
}

func TestShutdown(t *testing.T) {
	// Initialize first
	InitTelemetry(context.Background(), "test-service")
	assert.NotNil(t, tracerProvider)

	// Test shutdown
	Shutdown()

	// Note: Can't really test the internal state after shutdown
	// as the tracerProvider remains in memory but its internal state changes
}

func createTestCtx(method, target string, body []byte) (*fiber.Ctx, *fasthttp.RequestCtx) {
	app := fiber.New()

	// Buat fasthttp.RequestCtx manual
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	req.SetRequestURI(target)
	if body != nil {
		req.SetBody(body)
		req.Header.SetContentType("application/json")
	}

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	req.Header.CopyTo(&ctx.Request().Header)
	ctx.Request().SetBody(req.Body())

	return ctx, ctx.Context()
}

func TestStartSpanHandler(t *testing.T) {
	// Initialize first
	InitTelemetry(context.Background(), "test-service")

	// Create a test Fiber app and context
	c, _ := createTestCtx(fiber.MethodGet, "/test", nil)

	// Set a request ID in locals
	requestID := "test-request-id"
	c.Locals("requestid", requestID)

	// Call StartSpanHandler
	ctx, span := StartSpanHandler(c, "test-span")

	// Verify results
	assert.NotNil(t, ctx)
	assert.NotNil(t, span)

	// End the span
	span.End()
}

func TestStartSpan(t *testing.T) {
	// Initialize first
	InitTelemetry(context.Background(), "test-service")

	// Test with valid tracer
	ctx := context.Background()
	ctx, span := StartSpan(ctx, "test-span")

	assert.NotNil(t, ctx)
	assert.NotNil(t, span)
	assert.NotNil(t, span.Span)

	// End the span
	span.End()

	// Test with nil tracer
	// Temporarily set tracer to nil
	oldTracer := tracer
	tracer = nil

	ctx2, span2 := StartSpan(context.Background(), "test-span")
	assert.Nil(t, ctx2)
	assert.Nil(t, span2)

	// Restore tracer
	tracer = oldTracer
}

func TestFromContext(t *testing.T) {
	// Initialize first
	InitTelemetry(context.Background(), "test-service")

	// Create a context with a span
	ctx := context.Background()
	ctx, originalSpan := StartSpan(ctx, "test-span")

	// Test FromContext
	retrievedSpan := FromContext(ctx)

	assert.NotNil(t, retrievedSpan)
	assert.Equal(t, originalSpan.GetTraceID(), retrievedSpan.GetTraceID())

	// End the span
	originalSpan.End()
}

func TestToKeyValue(t *testing.T) {
	// Test various data types
	attributes := map[string]any{
		"string":  "test",
		"int":     42,
		"int64":   int64(42),
		"float64": 42.42,
		"bool":    true,
		"stringSlice": []string{
			"a",
			"b",
			"c",
		},
		"intSlice": []int{
			1,
			2,
			3,
		},
		"int64Slice": []int64{
			1,
			2,
			3,
		},
		"float64Slice": []float64{
			1.1,
			2.2,
			3.3,
		},
		"boolSlice": []bool{
			true,
			false,
			true,
		},
		"ignored": struct{}{}, // This type should be ignored
	}

	attrs := ToKeyValue(attributes)

	// Count how many attributes were converted (ignoring the struct type)
	assert.Equal(t, 10, len(attrs))

	// Check some specific values
	for _, attr := range attrs {
		switch attr.Key {
		case "string":
			assert.Equal(t, "test", attr.Value.AsString())
		case "int":
			assert.Equal(t, int64(42), attr.Value.AsInt64())
		case "bool":
			assert.Equal(t, true, attr.Value.AsBool())
		}
	}
}

func TestSpan_AddEventHelper(t *testing.T) {
	// Initialize first
	InitTelemetry(context.Background(), "test-service")

	// Create a mock span
	mockSpan := new(MockSpan)
	span := &Span{mockSpan}

	// Set up expectations
	attributes := map[string]any{
		"string": "test",
		"int":    42,
	}

	// The mock should expect AddEvent to be called with the name and attributes converted to KeyValue
	mockSpan.On("AddEvent", "test-event", mock.Anything).Return()

	// Call AddEventHelper
	span.AddEventHelper("test-event", attributes)

	// Verify expectations
	mockSpan.AssertExpectations(t)

	// Test with nil attributes
	span.AddEventHelper("test-event-nil", nil)
	mockSpan.AssertExpectations(t)
}

func TestSpan_RecordErrorHelper(t *testing.T) {
	// Initialize first
	InitTelemetry(context.Background(), "test-service")

	// Create a mock span
	mockSpan := new(MockSpan)
	span := &Span{mockSpan}

	// Set up expectations
	testErr := errors.New("test error")

	// The mock should expect RecordError and SetStatus to be called
	mockSpan.On("RecordError", testErr, mock.Anything).Return()
	mockSpan.On("SetStatus", codes.Error, "error message").Return()

	// Call RecordErrorHelper
	span.RecordErrorHelper(testErr, "error message")

	// Verify expectations
	mockSpan.AssertExpectations(t)

	// Test with custom stacktrace
	mockSpan.On("RecordError", testErr, mock.Anything).Return()
	mockSpan.On("SetStatus", codes.Error, "error with custom stack").Return()

	span.RecordErrorHelper(testErr, "error with custom stack")

	// Verify expectations
	mockSpan.AssertExpectations(t)
}

func TestSpan_GetTraceID(t *testing.T) {
	// Initialize first
	InitTelemetry(context.Background(), "test-service")

	_, span := StartSpan(context.Background(), "test-span")
	defer span.End()

	traceID := span.GetTraceID()

	// Validasi: traceID harus 32 karakter hex
	assert.Len(t, traceID, 32)
	assert.NotEqual(t, "00000000000000000000000000000000", traceID)
}

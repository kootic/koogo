package kootel

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	otellog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type OTelExporterType string

const (
	OTelExporterTypeConsole  OTelExporterType = "console"
	OTelExporterTypeOTLPgRPC OTelExporterType = "otlp-grpc"
	OTelExporterTypeNoop     OTelExporterType = "none"
)

type OTelConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	ExporterType   OTelExporterType
}

// InitializeOTel initializes the OpenTelemetry SDK and returns a function to shutdown the SDK.
// Other than the exporter protocol which is configured through the OTelConfig parameter,
// the SDK should be configured using standard OpenTelemetry environment variables.
// See: https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/
func InitializeOTel(ctx context.Context, cfg OTelConfig) (func(context.Context) error, error) {
	if cfg.ExporterType == OTelExporterTypeNoop {
		return func(ctx context.Context) error { return nil }, nil
	}

	res := newResource(cfg)

	lp, err := newLoggerProvider(ctx, res, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger provider: %w", err)
	}

	mp, err := newMeterProvider(ctx, res, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create meter provider: %w", err)
	}

	tp, err := newTracerProvider(ctx, res, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer provider: %w", err)
	}

	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	global.SetLoggerProvider(lp)

	shutdownFuncs := []func(context.Context) error{
		lp.Shutdown,
		mp.Shutdown,
		tp.Shutdown,
	}

	return func(ctx context.Context) error {
		var errs error

		for _, fn := range shutdownFuncs {
			if err := fn(ctx); err != nil {
				errs = errors.Join(errs, err)
			}
		}

		return nil
	}, nil
}

// newResource creates a new OTEL resource with the service name and version.
func newResource(cfg OTelConfig) *resource.Resource {
	hostName, _ := os.Hostname()

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.ServiceName),
		semconv.ServiceVersion(cfg.ServiceVersion),
		semconv.HostName(hostName),
		semconv.DeploymentEnvironment(cfg.Environment),
	)
}

// newLoggerProvider creates a new logger provider with the OTLP gRPC exporter.
func newLoggerProvider(ctx context.Context, res *resource.Resource, cfg OTelConfig) (*otellog.LoggerProvider, error) {
	var exporter otellog.Exporter

	var err error

	switch cfg.ExporterType {
	case OTelExporterTypeConsole:
		exporter, err = stdoutlog.New(stdoutlog.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout log exporter: %w", err)
		}
	case OTelExporterTypeOTLPgRPC:
		exporter, err = otlploggrpc.New(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
		}
	default:
		log.Println("Using no-op OpenTelemetry logger provider")
		return otellog.NewLoggerProvider(), nil
	}

	processor := otellog.NewBatchProcessor(exporter)

	lp := otellog.NewLoggerProvider(
		otellog.WithProcessor(processor),
		otellog.WithResource(res),
	)

	return lp, nil
}

// newMeterProvider creates a new meter provider with the OTLP gRPC exporter.
func newMeterProvider(ctx context.Context, res *resource.Resource, cfg OTelConfig) (*metric.MeterProvider, error) {
	var exporter metric.Exporter

	var err error

	switch cfg.ExporterType {
	case OTelExporterTypeConsole:
		exporter, err = stdoutmetric.New(stdoutmetric.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout metric exporter: %w", err)
		}
	case OTelExporterTypeOTLPgRPC:
		exporter, err = otlpmetricgrpc.New(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
		}
	default:
		log.Println("Using no-op OpenTelemetry meter provider")
		return metric.NewMeterProvider(), nil
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(res),
	)

	return mp, nil
}

// newTracerProvider creates a new tracer provider with the OTLP gRPC exporter.
func newTracerProvider(ctx context.Context, res *resource.Resource, cfg OTelConfig) (*trace.TracerProvider, error) {
	var exporter trace.SpanExporter

	var err error

	switch cfg.ExporterType {
	case OTelExporterTypeConsole:
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout trace exporter: %w", err)
		}
	case OTelExporterTypeOTLPgRPC:
		exporter, err = otlptracegrpc.New(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
		}
	default:
		log.Println("Using no-op OpenTelemetry tracer provider")
		return trace.NewTracerProvider(), nil
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	return tp, nil
}

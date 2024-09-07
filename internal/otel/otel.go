package otel

import (
	"context" // Provides context management capabilities, such as deadlines and cancellation.
	"errors"  // Standard package for creating and handling errors.
	"time"    // Standard package for time-related functions.

	// OpenTelemetry (OTel) packages for bridging, runtime instrumentation, and exporting telemetry data.
	"go.opentelemetry.io/contrib/bridges/otelslog"                      // Bridge for integrating OTel with standard logging.
	"go.opentelemetry.io/contrib/instrumentation/runtime"               // Provides runtime instrumentation for Go applications.
	"go.opentelemetry.io/otel"                                          // Core OpenTelemetry package.
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"       // HTTP exporter for sending logs to an OTel collector.
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp" // HTTP exporter for sending metrics to an OTel collector.
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"                 // Base package for tracing exporters.
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"   // HTTP exporter for sending traces to an OTel collector.
	"go.opentelemetry.io/otel/log/global"                               // Global logger management for OTel.
	"go.opentelemetry.io/otel/propagation"                              // Provides context propagation mechanisms.
	"go.opentelemetry.io/otel/sdk/log"                                  // SDK for logging with OTel.
	"go.opentelemetry.io/otel/sdk/metric"                               // SDK for metric collection and exporting with OTel.
	"go.opentelemetry.io/otel/sdk/trace"                                // SDK for trace management and exporting with OTel.
)

// schemaName is a constant defining the schema or standard URL used for instrumentation.
const schemaName = "https://github.com/grafana/docker-otel-lgtm"

// Global variables for tracing and logging.
// tracer is an OTel tracer used for creating spans (units of work).
var tracer = otel.Tracer(schemaName)

// logger is an OTel logger used for logging messages with the specified schema.
var logger = otelslog.NewLogger(schemaName)

// Setup initializes and configures the OpenTelemetry SDK for tracing, metrics, and logging.
// It returns a shutdown function for cleaning up resources and an error if something goes wrong.
func Setup(ctx context.Context) (shutdown func(context.Context) error, err error) {
	// shutdownFuncs holds cleanup functions for each OTel component (tracer, meter, logger).
	var shutdownFuncs []func(context.Context) error

	// shutdown function iterates over registered cleanup functions and invokes them.
	// It accumulates any errors encountered during the cleanup process.
	shutdown = func(ctx context.Context) error {
		var err error
		// Iterate over each cleanup function and invoke it.
		for _, fn := range shutdownFuncs {
			// Join errors from each shutdown function using errors.Join.
			err = errors.Join(err, fn(ctx))
		}
		// Clear the shutdownFuncs slice after invoking all cleanup functions.
		shutdownFuncs = nil
		// Return any accumulated error.
		return err
	}

	// handleErr is a helper function to handle errors during the setup process.
	// If an error occurs, it attempts to shutdown the OTel SDK and accumulates any additional errors.
	handleErr := func(inErr error) {
		// Join the incoming error with any errors from the shutdown process.
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up context propagation using a composite text map propagator.
	// This propagator includes both TraceContext and Baggage, which are common context formats.
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // Standard for passing trace information (trace IDs, span IDs) across services.
		propagation.Baggage{},      // Standard for passing arbitrary key-value pairs across services.
	)
	// Set the global text map propagator to the one just created.
	otel.SetTextMapPropagator(prop)

	// Create a new trace exporter to send trace data to an OTel collector over HTTP.
	traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient())
	if err != nil {
		// If the trace exporter fails to initialize, return the error immediately.
		return nil, err
	}

	// Create a new tracer provider that batches traces before sending them to the exporter.
	tracerProvider := trace.NewTracerProvider(trace.WithBatcher(traceExporter))
	if err != nil {
		// If tracer provider creation fails, handle the error by shutting down and returning.
		handleErr(err)
		return nil, err
	}
	// Register the tracer provider's shutdown function for later cleanup.
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	// Set the global tracer provider to the one just created.
	otel.SetTracerProvider(tracerProvider)

	// Create a new metric exporter to send metric data to an OTel collector over HTTP.
	metricExporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		// If the metric exporter fails to initialize, return the error immediately.
		return nil, err
	}

	// Create a new meter provider that periodically reads and exports metrics using the metric exporter.
	meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(metricExporter)))
	if err != nil {
		// If meter provider creation fails, handle the error by shutting down and returning.
		handleErr(err)
		return nil, err
	}
	// Register the meter provider's shutdown function for later cleanup.
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	// Set the global meter provider to the one just created.
	otel.SetMeterProvider(meterProvider)

	// Create a new log exporter to send log data to an OTel collector over HTTP.
	// The WithInsecure option is used to allow insecure HTTP connections.
	logExporter, err := otlploghttp.New(ctx, otlploghttp.WithInsecure())
	if err != nil {
		// If the log exporter fails to initialize, return the error immediately.
		return nil, err
	}

	// Create a new logger provider that batches logs before sending them to the exporter.
	loggerProvider := log.NewLoggerProvider(log.WithProcessor(log.NewBatchProcessor(logExporter)))
	if err != nil {
		// If logger provider creation fails, handle the error by shutting down and returning.
		handleErr(err)
		return nil, err
	}
	// Register the logger provider's shutdown function for later cleanup.
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	// Set the global logger provider to the one just created.
	global.SetLoggerProvider(loggerProvider)

	// Start runtime instrumentation to collect runtime metrics (e.g., memory stats) with a minimum read interval of 1 second.
	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		// If runtime instrumentation fails, log the error using the OTel logger.
		logger.ErrorContext(ctx, "otel runtime instrumentation failed:", "error", err)
	}

	// Return the shutdown function and any error encountered during setup (nil if none).
	return shutdown, nil
}

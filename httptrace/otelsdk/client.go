package otelsdk

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"google.golang.org/grpc"
)

func StartClient(ctx context.Context) (func(context.Context) error, error) {
	res, err := newResource(ctx)
	if err != nil {
		return nil, err
	}

	traceExporter, err := newExporter(ctx)
	if err != nil {
		return nil, err
	}

	idg := xray.NewIDGenerator()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	return func(context.Context) (err error) {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		err = tp.Shutdown(ctx)
		if err != nil {
			return err
		}
		return nil
	}, nil
}

func newResource(ctx context.Context) (*resource.Resource, error) {
	resourceAttributes := []attribute.KeyValue{
		semconv.ServiceName("adot-tracing-sample"),
		semconv.ServiceVersion("1.0.0"),
	}

	ecsResourceDetector := ecs.NewResourceDetector()
	ecsRes, err := ecsResourceDetector.Detect(ctx)
	if err != nil {
		return nil, err
	}
	if ecsRes.Attributes() != nil {
		resourceAttributes = append(resourceAttributes, ecsRes.Attributes()...)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		resourceAttributes...,
	)

	return res, nil
}

func newExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "0.0.0.0:4317"
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithDialOption(grpc.WithBlock()))
	if err != nil {
		return nil, err
	}

	return traceExporter, nil
}

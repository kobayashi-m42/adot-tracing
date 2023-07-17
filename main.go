package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/kobayashi-m42/adot-tracing/httptrace/otelsdk"
	"github.com/kobayashi-m42/adot-tracing/httptrace/xraysdk"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx := context.Background()

	shutdown, err := otelsdk.StartClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(ctx)

	otelHttpClient := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	s3ClientForOtelSdk, err := otelsdk.NewS3Client(ctx)
	s3ClientForXraySdk, err := xraysdk.NewS3Client(ctx)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("healthcheck"))
	}))
	mux.Handle("/xray-sdk/outgoing-http-call",
		xray.Handler(
			xray.NewFixedSegmentNamer("/xray-sdk/outgoing-http-call"),
			http.HandlerFunc(xraysdk.OutgoingHttpCall),
		),
	)
	mux.Handle(
		"/xray-sdk/aws-sdk-call",
		xray.Handler(
			xray.NewFixedSegmentNamer("/xray-sdk/aws-sdk-call"),
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					xraysdk.AwsSdkCall(w, r, s3ClientForXraySdk)
				},
			),
		),
	)
	mux.Handle("/otel-sdk/outgoing-http-call",
		otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			otelsdk.OutgoingHttpCall(w, r, otelHttpClient)
		}),
			"server",
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		),
	)
	mux.Handle("/otel-sdk/aws-sdk-call",
		otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			otelsdk.AwsSdkCall(w, r, s3ClientForOtelSdk)
		}),
			"server",
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		),
	)
	server := &http.Server{
		Addr:    ":80",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

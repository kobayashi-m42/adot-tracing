package otelsdk

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type response struct {
	TraceID string `json:"traceId"`
}

var tracer = otel.Tracer("github.com/kobayashi-m42/adot-tracing/httptrace/otesdk")

func AwsSdkCall(w http.ResponseWriter, r *http.Request, s3 *S3Client) {
	ctx, span := tracer.Start(
		r.Context(),
		"aws-sdk-call",
	)
	defer span.End()

	if _, err := s3.Client.ListBuckets(ctx, nil); err != nil {
		log.Println(err)
	}

	writeResponse(span, w)
}

func OutgoingHttpCall(w http.ResponseWriter, r *http.Request, client http.Client) {
	ctx, span := tracer.Start(
		r.Context(),
		"outgoing-http-call",
	)
	defer span.End()

	req, _ := http.NewRequestWithContext(ctx, "GET", "https://aws.amazon.com/", nil)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	writeResponse(span, w)
}

func getXrayTraceID(span trace.Span) string {
	xrayTraceID := span.SpanContext().TraceID().String()
	return fmt.Sprintf("1-%s-%s", xrayTraceID[0:8], xrayTraceID[8:])
}

func writeResponse(span trace.Span, w http.ResponseWriter) {
	xrayTraceID := getXrayTraceID(span)
	payload, _ := json.Marshal(response{TraceID: xrayTraceID})
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

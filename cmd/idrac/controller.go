package main

import (
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

func initTracer() {
	// Create jaeger exporter to be able to retrieve
	// the collected spans.
	jExporter, err := jaeger.NewRawExporter(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")),
	)

	if err != nil {
		log.Fatal(err)
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp := sdktrace.NewTracerProvider(
		// set Span Name
		sdktrace.WithResource(resource.NewWithAttributes(semconv.ServiceNameKey.String("Alochym"))),
		sdktrace.WithSyncer(jExporter),
	)
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

func main() {
	initTracer()

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "Hello")

	// http.HandleFunc("/hello", helloHandler)
	http.Handle	("/hello", otelHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	// get global Trace Provider
	tr := otel.GetTracerProvider().Tracer("Alochym")

	// get context from http request
	ctx := req.Context()

	// start tracing
	_, sp := tr.Start(ctx, "helloHandler")

	// close tracing
	defer sp.End()

	// http response
	_, _ = io.WriteString(w, "Hello, world!\n")
}

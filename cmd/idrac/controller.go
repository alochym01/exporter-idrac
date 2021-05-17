package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

		// register jaeger exporter
		sdktrace.WithSyncer(jExporter),
	)

	// handle err
	if err != nil {
		log.Fatal(err)
	}

	// set global opentelemetry
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

func main() {
	initTracer()

	// using auto-instrument opentelemetry
	http.Handle("/hello", otelhttp.NewHandler(http.HandlerFunc(helloHandler), "root"))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	// get context from http request
	ctx := req.Context()

	// passing span context to other func
	usingCtx(ctx)

	// http response
	_, _ = io.WriteString(w, "Hello, world!\n")
}

func usingCtx(ctx context.Context) {
	tr := otel.GetTracerProvider().Tracer("Alochym")

	// start tracing
	_, sp := tr.Start(ctx, "usingCtx")

	// set attribute
	sp.SetAttributes(
		attribute.String("usingCtx Function", "OLACHYM"),
		attribute.String("Function", "OLACHYM 02"),
	)

	// close tracing
	defer sp.End()

	// doing something
	fmt.Println("alochym")
}

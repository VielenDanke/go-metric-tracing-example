package trace

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func HTTPHandler(handler http.Handler, name string) http.Handler {
	return otelhttp.NewHandler(handler, name, otelhttp.WithTracerProvider(otel.GetTracerProvider()))
}

func HTTPHandlerFunc(handler http.HandlerFunc, name string) http.HandlerFunc {
	return otelhttp.NewHandler(handler, name, otelhttp.WithTracerProvider(otel.GetTracerProvider())).ServeHTTP
}

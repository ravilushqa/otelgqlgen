module github.com/ravilushqa/otelgqlgen/example

go 1.16

replace github.com/ravilushqa/otelgqlgen => ../

require (
	github.com/99designs/gqlgen v0.14.0
	github.com/ravilushqa/otelgqlgen v0.3.0
	github.com/vektah/gqlparser/v2 v2.2.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.3.0
	go.opentelemetry.io/otel/sdk v1.3.0
	go.opentelemetry.io/otel/trace v1.3.0
)

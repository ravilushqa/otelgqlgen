module github.com/ravilushqa/otelgqlgen/test

go 1.15

require (
	github.com/99designs/gqlgen v0.14.0
	github.com/ravilushqa/otelgqlgen v0.23.0
	github.com/stretchr/testify v1.7.0
	github.com/vektah/gqlparser/v2 v2.2.0
	go.opentelemetry.io/otel v1.0.0-RC3
	go.opentelemetry.io/otel/sdk v1.0.0-RC3
	go.opentelemetry.io/otel/trace v1.0.0-RC3
)

replace github.com/ravilushqa/otelgqlgen => ../

module github.com/ravilushqa/otelgqlgen/example

go 1.16

replace github.com/ravilushqa/otelgqlgen => ../

require (
	github.com/99designs/gqlgen v0.16.0
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/ravilushqa/otelgqlgen v0.4.1
	github.com/vektah/gqlparser/v2 v2.3.1
	go.opentelemetry.io/otel v1.4.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.4.0
	go.opentelemetry.io/otel/sdk v1.4.0
	go.opentelemetry.io/otel/trace v1.4.0
	golang.org/x/sys v0.0.0-20220204135822-1c1b9b1eba6a // indirect
)

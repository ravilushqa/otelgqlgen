module github.com/ravilushqa/otelgqlgen/example

go 1.16

replace github.com/ravilushqa/otelgqlgen => ../

require (
	github.com/99designs/gqlgen v0.17.2
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/ravilushqa/otelgqlgen v0.6.0
	github.com/vektah/gqlparser/v2 v2.4.1
	go.opentelemetry.io/otel v1.6.3
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.6.3
	go.opentelemetry.io/otel/sdk v1.6.3
	go.opentelemetry.io/otel/trace v1.6.3
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
)

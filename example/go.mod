module github.com/ravilushqa/otelgqlgen/example

go 1.19

replace github.com/ravilushqa/otelgqlgen => ../

require (
	github.com/99designs/gqlgen v0.17.20
	github.com/ravilushqa/otelgqlgen v0.0.0
	github.com/vektah/gqlparser/v2 v2.5.1
	go.opentelemetry.io/otel v1.15.1
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.9.0
	go.opentelemetry.io/otel/sdk v1.15.1
	go.opentelemetry.io/otel/trace v1.15.1
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	go.opentelemetry.io/contrib v1.16.1 // indirect
	golang.org/x/sys v0.7.0 // indirect
)

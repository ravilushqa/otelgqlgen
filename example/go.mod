module github.com/ravilushqa/otelgqlgen/example

go 1.20

replace github.com/ravilushqa/otelgqlgen => ../

require (
	github.com/99designs/gqlgen v0.17.41
	github.com/ravilushqa/otelgqlgen v0.0.0
	github.com/vektah/gqlparser/v2 v2.5.10
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.9.0
	go.opentelemetry.io/otel/sdk v1.21.0
	go.opentelemetry.io/otel/trace v1.21.0
)

require (
	github.com/google/uuid v1.5.0 // indirect
	github.com/sosodev/duration v1.2.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.3 // indirect; indirect // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	go.opentelemetry.io/contrib v1.21.1 // indirect
	golang.org/x/sys v0.15.0 // indirect
)

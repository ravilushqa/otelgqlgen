module github.com/ravilushqa/otelgqlgen/example

go 1.23.8

toolchain go1.23.9

replace github.com/ravilushqa/otelgqlgen => ../

require (
	github.com/99designs/gqlgen v0.17.73
	github.com/ravilushqa/otelgqlgen v0.17.0
	github.com/vektah/gqlparser/v2 v2.5.27
	go.opentelemetry.io/otel v1.36.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.36.0
	go.opentelemetry.io/otel/sdk v1.36.0
	go.opentelemetry.io/otel/trace v1.36.0
)

require (
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
)

require (
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect; indirect // indirect
	go.opentelemetry.io/contrib v1.36.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
)

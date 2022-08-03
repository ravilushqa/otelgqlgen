# OpenTelemetry-Go gqlgen Instrumentation

[![ci](https://github.com/ravilushqa/otelgqlgen/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/ravilushqa/otelgqlgen/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ravilushqa/otelgqlgen)](https://goreportcard.com/report/github.com/ravilushqa/otelgqlgen)
[![codecov](https://codecov.io/gh/ravilushqa/otelgqlgen/branch/main/graph/badge.svg?token=4LXALH3K72)](https://codecov.io/gh/ravilushqa/otelgqlgen)
[![Go Reference](https://pkg.go.dev/badge/github.com/ravilushqa/otelgqlgen.svg)](https://pkg.go.dev/github.com/ravilushqa/otelgqlgen)

It is an OpenTelemetry instrumentation for Golang 99designs/gqlgen, a port from https://github.com/open-telemetry/opentelemetry-go-contrib/pull/761.

It can only instrument traces for the present.

## Installation

```
$ go get github.com/ravilushqa/otelgqlgen
```

## Example

See [./example](./example).

## Compatibility

This project is tested on the following systems.

| OS      | Go Version | Architecture |
| ------- | ---------- | ------------ |
| Ubuntu  | 1.19       | amd64        |
| Ubuntu  | 1.18       | amd64        |
| Ubuntu  | 1.19       | 386          |
| Ubuntu  | 1.18       | 386          |
| MacOS   | 1.19       | amd64        |
| MacOS   | 1.18       | amd64        |
| Windows | 1.19       | amd64        |
| Windows | 1.18       | amd64        |
| Windows | 1.19       | 386          |
| Windows | 1.18       | 386          |

While this project should work for other systems, no compatibility guarantees
are made for those systems currently.

The project follows the [Release Policy](https://golang.org/doc/devel/release#policy) to support major Go releases.


## References

- [GraphQL](https://graphql.org/)
- [gqlgen](https://gqlgen.com)

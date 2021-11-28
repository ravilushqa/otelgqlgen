# OpenTelemetry-Go gqlgen Instrumentation

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
| Ubuntu  | 1.17       | amd64        |
| Ubuntu  | 1.16       | amd64        |
| Ubuntu  | 1.17       | 386          |
| Ubuntu  | 1.16       | 386          |
| MacOS   | 1.17       | amd64        |
| MacOS   | 1.16       | amd64        |
| Windows | 1.17       | amd64        |
| Windows | 1.16       | amd64        |
| Windows | 1.17       | 386          |
| Windows | 1.16       | 386          |

While this project should work for other systems, no compatibility guarantees
are made for those systems currently.

The project follows the [Release Policy](https://golang.org/doc/devel/release#policy) to support major Go releases.


## References

- [GraphQL](https://graphql.org/)
- [gqlgen](https://gqlgen.com)

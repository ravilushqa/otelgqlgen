// Copyright Ravil Galaktionov
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otelgqlgen

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"

	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	tracerName      = "github.com/ravilushqa/otelgqlgen"
	extensionName   = "OpenTelemetry"
	complexityLimit = "ComplexityLimit"
)

// Tracer is a GraphQL extension that traces GraphQL requests.
type Tracer struct {
	complexityExtensionName     string
	tracer                      oteltrace.Tracer
	requestVariablesBuilderFunc RequestVariablesBuilderFunc
	shouldCreateSpanFromFields  FieldsPredicateFunc
	spanKindSelector            SpanKindSelectorFunc
	setHTTPStatusCode           bool
}

var _ interface {
	graphql.HandlerExtension
	graphql.ResponseInterceptor
	graphql.FieldInterceptor
} = Tracer{}

// ExtensionName returns the extension name.
func (gqlTracer Tracer) ExtensionName() string {
	return extensionName
}

// Validate checks if the extension is configured properly.
func (gqlTracer Tracer) Validate(_ graphql.ExecutableSchema) error {
	return nil
}

// InterceptResponse intercepts the incoming request.
func (gqlTracer Tracer) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	if !graphql.HasOperationContext(ctx) {
		return next(ctx)
	}

	opName := operationName(ctx)
	spanKind := gqlTracer.spanKindSelector(opName)
	ctx, span := gqlTracer.tracer.Start(ctx, opName, oteltrace.WithSpanKind(spanKind))
	defer span.End()
	if !span.IsRecording() {
		return next(ctx)
	}

	oc := graphql.GetOperationContext(ctx)

	span.SetAttributes(
		RequestQuery(oc.RawQuery),
	)
	complexityExtension := gqlTracer.complexityExtensionName
	if complexityExtension == "" {
		complexityExtension = complexityLimit
	}
	complexityStats, ok := oc.Stats.GetExtension(complexityExtension).(*extension.ComplexityStats)
	if !ok {
		// complexity extension is not used
		complexityStats = &extension.ComplexityStats{}
	}

	if complexityStats.ComplexityLimit > 0 {
		span.SetAttributes(
			RequestComplexityLimit(int64(complexityStats.ComplexityLimit)),
			RequestOperationComplexity(int64(complexityStats.Complexity)),
		)
	}

	if gqlTracer.requestVariablesBuilderFunc != nil {
		span.SetAttributes(gqlTracer.requestVariablesBuilderFunc(oc.Variables)...)
	}

	resp := next(ctx)
	if gqlTracer.setHTTPStatusCode {
		var statusCode = http.StatusOK
		if resp != nil && len(resp.Errors) > 0 {
			statusCode = http.StatusInternalServerError
			span.SetStatus(codes.Error, resp.Errors.Error())
			span.RecordError(fmt.Errorf("graphql response errors: %v", resp.Errors.Error()))
			span.SetAttributes(ResolverErrors(resp.Errors)...)
			if resp.Errors[0].Extensions["code"] != nil {
				var code = resp.Errors[0].Extensions["code"]
				if reflect.TypeOf(code).Kind() == reflect.Int {
					statusCode = code.(int)
				} else if reflect.TypeOf(code).Kind() == reflect.Int32 {
					statusCode = int(code.(int32))
				} else if reflect.TypeOf(code).Kind() == reflect.String {
					status, err := strconv.Atoi(code.(string))
					if err == nil {
						statusCode = status
					}
				}
			}
		} else {
			span.SetStatus(codes.Ok, "Finished successfully")
		}

		// Stable: https://opentelemetry.io/docs/specs/semconv/http/http-spans/#common-attributes
		span.SetAttributes(semconv.HTTPResponseStatusCode(statusCode))
		// Experimental: https://github.com/open-telemetry/opentelemetry-specification/blob/v1.20.0/specification/trace/semantic_conventions/http.md#common-attributes
		span.SetAttributes(semconv.HTTPStatusCode(statusCode)) // nolint:staticcheck
	}

	return resp
}

// InterceptField intercepts the incoming request.
func (gqlTracer Tracer) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)
	if !gqlTracer.shouldCreateSpanFromFields(fc) {
		return next(ctx)
	}
	name := fc.Field.ObjectDefinition.Name + "/" + fc.Field.Name
	spanKind := gqlTracer.spanKindSelector(name)
	ctx, span := gqlTracer.tracer.Start(ctx,
		name,
		oteltrace.WithSpanKind(spanKind),
	)
	defer span.End()
	if !span.IsRecording() {
		return next(ctx)
	}

	span.SetAttributes(
		ResolverPath(fc.Path().String()),
		ResolverObject(fc.Field.ObjectDefinition.Name),
		ResolverField(fc.Field.Name),
		ResolverAlias(fc.Field.Alias),
	)
	span.SetAttributes(ResolverArgs(fc.Field.Arguments)...)

	resp, err := next(ctx)

	errList := graphql.GetFieldErrors(ctx, fc)
	var statusCode = http.StatusOK
	if len(errList) != 0 {
		statusCode = http.StatusInternalServerError
		span.SetStatus(codes.Error, errList.Error())
		span.RecordError(fmt.Errorf("graphql field errors: %v", errList.Error()))
		span.SetAttributes(ResolverErrors(errList)...)
		if errList[0].Extensions["code"] != nil {
			var code = errList[0].Extensions["code"]
			if reflect.TypeOf(code).Kind() == reflect.Int {
				statusCode = code.(int)
			} else if reflect.TypeOf(code).Kind() == reflect.Int32 {
				statusCode = int(code.(int32))
			} else if reflect.TypeOf(code).Kind() == reflect.String {
				status, err := strconv.Atoi(code.(string))
				if err == nil {
					statusCode = status
				}
			}
		}
	} else {
		span.SetStatus(codes.Ok, "Finished successfully")
	}

	// Stable: https://opentelemetry.io/docs/specs/semconv/http/http-spans/#common-attributes
	span.SetAttributes(semconv.HTTPResponseStatusCode(statusCode))
	// Experimental: https://github.com/open-telemetry/opentelemetry-specification/blob/v1.20.0/specification/trace/semantic_conventions/http.md#common-attributes
	span.SetAttributes(semconv.HTTPStatusCode(statusCode)) // nolint:staticcheck

	return resp, err
}

// Middleware sets up a handler to start tracing the incoming
// requests.  The service parameter should describe the name of the
// (virtual) server handling the request. extension parameter may be empty string.
func Middleware(opts ...Option) Tracer {
	cfg := config{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	if cfg.RequestVariablesBuilder == nil {
		cfg.RequestVariablesBuilder = RequestVariables
	}
	if cfg.ShouldCreateSpanFromFields == nil {
		cfg.ShouldCreateSpanFromFields = alwaysTrue()
	}
	if cfg.SpanKindSelectorFunc == nil {
		cfg.SpanKindSelectorFunc = alwaysServer()
	}

	tracer := cfg.TracerProvider.Tracer(
		tracerName,
		oteltrace.WithInstrumentationVersion(otelcontrib.Version()),
	)

	return Tracer{
		tracer:                      tracer,
		requestVariablesBuilderFunc: cfg.RequestVariablesBuilder,
		shouldCreateSpanFromFields:  cfg.ShouldCreateSpanFromFields,
		spanKindSelector:            cfg.SpanKindSelectorFunc,
		setHTTPStatusCode:           cfg.SetHTTPStatusCode,
	}

}

// alwaysTrue returns a FieldsPredicateFunc that always returns true.
func alwaysTrue() FieldsPredicateFunc {
	return func(_ *graphql.FieldContext) bool {
		return true
	}
}

func alwaysServer() SpanKindSelectorFunc {
	return func(_ string) oteltrace.SpanKind {
		return oteltrace.SpanKindServer
	}
}

func operationName(ctx context.Context) string {
	opContext := graphql.GetOperationContext(ctx)
	if opName := opContext.OperationName; opName != "" {
		return opName
	}
	if opContext.Operation != nil && opContext.Operation.Name != "" {
		return opContext.Operation.Name
	}
	return GetOperationName(ctx)
}

type operationNameCtxKey struct{}

// SetOperationName adds the operation name to the context so that the interceptors can use it.
// It will replace the operation name if it already exists in the context.
// example:
//
//		ctx = otelgqlgen.SetOperationName(r.Context(), "my-operation")
//	 	r = r.WithContext(ctx)
func SetOperationName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, operationNameCtxKey{}, name)
}

// GetOperationName gets the operation name from the context.
func GetOperationName(ctx context.Context) string {
	if oc, _ := ctx.Value(operationNameCtxKey{}).(string); oc != "" {
		return oc
	}
	return "nameless-operation"
}

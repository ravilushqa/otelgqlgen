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

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"

	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	tracerName      = "github.com/ravilushqa/otelgqlgen"
	extensionName   = "OpenTelemetry"
	complexityLimit = "ComplexityLimit"
)

type Tracer struct {
	complexityExtensionName     string
	tracer                      oteltrace.Tracer
	requestVariablesBuilderFunc RequestVariablesBuilderFunc
}

var _ interface {
	graphql.HandlerExtension
	graphql.ResponseInterceptor
	graphql.FieldInterceptor
} = Tracer{}

func (a Tracer) ExtensionName() string {
	return extensionName
}

func (a Tracer) Validate(_ graphql.ExecutableSchema) error {
	return nil
}

func (a Tracer) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	ctx, span := a.tracer.Start(ctx, operationName(ctx), oteltrace.WithSpanKind(oteltrace.SpanKindServer))
	defer span.End()
	if !span.IsRecording() {
		return next(ctx)
	}

	oc := graphql.GetOperationContext(ctx)

	span.SetAttributes(
		RequestQuery(oc.RawQuery),
	)
	complexityExtension := a.complexityExtensionName
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

	if a.requestVariablesBuilderFunc != nil {
		span.SetAttributes(a.requestVariablesBuilderFunc(oc.Variables)...)
	}

	resp := next(ctx)
	if resp != nil && len(resp.Errors) > 0 {
		span.SetStatus(codes.Error, resp.Errors.Error())
		span.RecordError(fmt.Errorf(resp.Errors.Error()))
		span.SetAttributes(ResolverErrors(resp.Errors)...)
	}

	return resp
}

func (a Tracer) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)
	ctx, span := a.tracer.Start(ctx,
		fc.Field.ObjectDefinition.Name+"/"+fc.Field.Name,
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
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
	if len(errList) != 0 {
		span.SetStatus(codes.Error, errList.Error())
		span.RecordError(fmt.Errorf(errList.Error()))
		span.SetAttributes(ResolverErrors(errList)...)
	}

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

	tracer := cfg.TracerProvider.Tracer(
		tracerName,
		oteltrace.WithInstrumentationVersion(otelcontrib.SemVersion()),
	)

	return Tracer{
		tracer:                      tracer,
		requestVariablesBuilderFunc: cfg.RequestVariablesBuilder,
	}

}

func operationName(ctx context.Context) string {
	requestContext := graphql.GetOperationContext(ctx)
	if requestContext.Doc != nil && len(requestContext.Doc.Operations) != 0 {
		op := requestContext.Doc.Operations[0]
		if op.Name != "" {
			return op.Name
		}
	}
	return GetOperationName(ctx)
}

type operationNameCtxKey struct{}

// SetOperationName adds the operation name to the context so that the interceptors can use it.
// It will replace the operation name if it already exists in the context.
// example:
// 	ctx = otelgqlgen.SetOperationName(r.Context(), "my-operation")
//  	r = r.WithContext(ctx)
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

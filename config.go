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
	"github.com/99designs/gqlgen/graphql"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type FieldsPredicateFunc func(ctx *graphql.FieldContext) bool

type SpanKindSelectorFunc func(ctx *graphql.FieldContext) trace.SpanKind

// config is used to configure the mongo tracer.
type config struct {
	TracerProvider             trace.TracerProvider
	Tracer                     trace.Tracer
	ComplexityExtensionName    string
	RequestVariablesBuilder    RequestVariablesBuilderFunc
	ShouldCreateSpanFromFields FieldsPredicateFunc
	SpanKindSelectorFunc       SpanKindSelectorFunc
}

// RequestVariablesBuilderFunc is the signature of the function
// used to build the request variables attributes.
type RequestVariablesBuilderFunc func(requestVariables map[string]interface{}) []attribute.KeyValue

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider trace.TracerProvider) Option {
	return optionFunc(func(cfg *config) {
		cfg.TracerProvider = provider
	})
}

// WithComplexityExtensionName specifies complexity extension name.
func WithComplexityExtensionName(complexityExtensionName string) Option {
	return optionFunc(func(cfg *config) {
		cfg.ComplexityExtensionName = complexityExtensionName
	})
}

// WithRequestVariablesAttributesBuilder allows specifying a custom function
// to handle the building of the attributes for the variables.
func WithRequestVariablesAttributesBuilder(builder RequestVariablesBuilderFunc) Option {
	return optionFunc(func(cfg *config) {
		cfg.RequestVariablesBuilder = builder
	})
}

// WithoutVariables allows disabling the variables attributes.
func WithoutVariables() Option {
	return optionFunc(func(cfg *config) {
		cfg.RequestVariablesBuilder = func(requestVariables map[string]interface{}) []attribute.KeyValue {
			return nil
		}
	})
}

// WithCreateSpanFromFields allows specifying a custom function
// to handle the creation or not of spans regarding the GraphQL context fields.
func WithCreateSpanFromFields(predicate FieldsPredicateFunc) Option {
	return optionFunc(func(cfg *config) {
		cfg.ShouldCreateSpanFromFields = predicate
	})
}

// WithSpanKindSelector allows specifying a custom function that defines the SpanKind based on the name
func WithSpanKindSelector(spanKindSelector SpanKindSelectorFunc) Option {
	return optionFunc(func(cfg *config) {
		cfg.SpanKindSelectorFunc = spanKindSelector
	})
}

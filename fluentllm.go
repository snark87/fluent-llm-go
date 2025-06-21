// Package fluentllm provides a fluent interface for working with Large Language Models.
package fluentllm

import (
	"context"
	"encoding/json"
)

type Model interface {
	Prompt(string) PromptBuilder
	NewSchema() SchemaBuilder
}

type PromptBuilder interface {
	Execute(context.Context) (LLMResponse, error)
	WithSchema(Schema) PromptBuilder
}

type LLMResponse interface {
	Text() string
}

func AsStructuredResponse[T any](response LLMResponse) (T, error) {
	var output T
	err := json.Unmarshal([]byte(response.Text()), &output)
	if err != nil {
		return output, err
	}

	return output, nil
}

type Schema interface {
	BuildArg
}

type SchemaBuilder interface {
	FromJSONSchema(schema []byte) (Schema, error)
	FromGoValue(value any) (Schema, error)

	Object(name string, buildArgs ...BuildArg) Schema
	Str(name string, buildArgs ...BuildArg) Schema
	Int(name string, buildArgs ...BuildArg) Schema

	Description(description string) BuildArg
	Required() BuildArg
}

type BuildArg interface {
	BuildArgType() BuildArgType
}

type BuildArgType string

const (
	BuildArgTypeSchema BuildArgType = "schema"
	BuildArgTypeOption BuildArgType = "option"
)

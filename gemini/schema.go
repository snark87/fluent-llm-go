package gemini

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/snark87/fluentllm"
	"google.golang.org/genai"
)

type SchemaBuilder struct {
}

var _ fluentllm.SchemaBuilder = (*SchemaBuilder)(nil)

type Schema struct {
	name     string
	required bool
	schema   *genai.Schema
}

type schemaBuildOption interface {
	fluentllm.BuildArg

	apply(schema *Schema)
}

var _ fluentllm.Schema = (*Schema)(nil)

func (s *Schema) applyBuildOption(arg fluentllm.BuildArg) {
	switch arg.BuildArgType() {
	case fluentllm.BuildArgTypeOption:
		buildArg, ok := arg.(schemaBuildOption)
		if !ok {
			panic(fmt.Sprintf("buildArg %+v is not a valid Gemini BuildOption", arg))
		}
		buildArg.apply(s)
	default:
		panic(fmt.Sprintf("buildArg of type %s is not supported in Gemini Schema", arg.BuildArgType()))
	}
}

func (sc *SchemaBuilder) FromJSONSchema(jsonSchema []byte) (fluentllm.Schema, error) {
	var schema genai.Schema
	if err := json.Unmarshal(jsonSchema, &schema); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON schema: %w", err)
	}
	return &Schema{
		name:   schema.Title,
		schema: &schema,
	}, nil
}

func (sc *SchemaBuilder) FromGoValue(value any) (fluentllm.Schema, error) {
	reflector := jsonschema.Reflector{
		DoNotReference: true,
	}
	schema := reflector.Reflect(value)
	bytes, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Go value to JSON schema: %w", err)
	}
	return sc.FromJSONSchema(bytes)
}

func (sc *SchemaBuilder) Object(name string, buildArgs ...fluentllm.BuildArg) fluentllm.Schema {
	schema := &Schema{
		name: name,
		schema: &genai.Schema{
			Type:       genai.TypeObject,
			Title:      name,
			Properties: make(map[string]*genai.Schema),
			Required:   make([]string, 0),
		},
	}

	for _, arg := range buildArgs {
		switch arg.BuildArgType() {
		case fluentllm.BuildArgTypeSchema:
			schemaArg, ok := arg.(*Schema)
			if !ok {
				panic(fmt.Sprintf("buildArg %+v is not a valid Gemini Schema", arg))
			}
			schema.schema.Properties[schemaArg.name] = schemaArg.schema
			if schemaArg.required {
				schema.schema.Required = append(schema.schema.Required, schemaArg.name)
			}

		case fluentllm.BuildArgTypeOption:
			schema.applyBuildOption(arg)
		default:
			panic(fmt.Sprintf("buildArg of type %s is not supported in Gemini Schema", arg.BuildArgType()))
		}
	}

	return schema
}

func (sc *SchemaBuilder) Str(name string, buildArgs ...fluentllm.BuildArg) fluentllm.Schema {
	schema := &Schema{
		name: name,
		schema: &genai.Schema{
			Type:  genai.TypeString,
			Title: name,
		},
	}
	for _, arg := range buildArgs {
		schema.applyBuildOption(arg)
	}
	return schema
}

func (sc *SchemaBuilder) Int(name string, buildArgs ...fluentllm.BuildArg) fluentllm.Schema {
	schema := &Schema{
		name: name,
		schema: &genai.Schema{
			Type:  genai.TypeInteger,
			Title: name,
		},
	}

	return schema
}

type descriptionOption struct {
	description string
}

func (d *descriptionOption) BuildArgType() fluentllm.BuildArgType {
	return fluentllm.BuildArgTypeOption
}

func (d *descriptionOption) apply(schema *Schema) {
	schema.schema.Description = d.description
}

func (sc *SchemaBuilder) Description(description string) fluentllm.BuildArg {
	return &descriptionOption{description: description}
}

type requiredOption struct{}

func (r *requiredOption) BuildArgType() fluentllm.BuildArgType {
	return fluentllm.BuildArgTypeOption
}

func (r *requiredOption) apply(schema *Schema) {
	schema.required = true
}

func (sc *SchemaBuilder) Required() fluentllm.BuildArg {
	return &requiredOption{}
}

func (gs *Schema) BuildArgType() fluentllm.BuildArgType {
	return fluentllm.BuildArgTypeSchema
}

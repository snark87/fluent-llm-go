package gemini

import (
	"context"
	"fmt"

	"google.golang.org/genai"

	"github.com/snark87/fluentllm"
	"github.com/snark87/fluentllm/internal/lazy"
)

type Model struct {
	modelName   string
	credentials Credentials

	client lazy.LazyWithError[*genai.Client]
}

var _ fluentllm.Model = (*Model)(nil)

type ModelOption interface {
	ApplyToModel(*Model)
}

func NewModel(modelName string, opts ...ModelOption) *Model {
	model := &Model{
		modelName: modelName,
	}

	for _, opt := range opts {
		opt.ApplyToModel(model)
	}

	return model
}

type PromptBuilder struct {
	model        *Model
	simplePrompt string
	schema       *genai.Schema
}

type GeminiResponse struct {
	model    *Model
	text     lazy.Lazy[string]
	response *genai.GenerateContentResponse
}

func (m *Model) Prompt(prompt string) fluentllm.PromptBuilder {
	return &PromptBuilder{
		model:        m,
		simplePrompt: prompt,
	}
}

func (m *Model) NewSchema() fluentllm.SchemaBuilder {
	return &SchemaBuilder{}
}

func (m *Model) String() string {
	return fmt.Sprintf("(%s)", m.modelName)
}

func (m *Model) setCredentials(credentials *Credentials) {
	if credentials.APIKey != m.credentials.APIKey { // pragma: allowlist secret
		m.client = lazy.NewWithError(m.newClient)
	}
	m.credentials = *credentials
}

func (m *Model) newClient(ctx context.Context) (*genai.Client, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: m.credentials.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return client, nil
}

func (p *PromptBuilder) Execute(ctx context.Context) (fluentllm.LLMResponse, error) {
	client, err := p.model.client.Get(ctx)
	if err != nil {
		return nil, err
	}
	var config *genai.GenerateContentConfig
	if p.schema != nil {
		config = &genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema:   p.schema,
		}
	}
	response, err := client.Models.GenerateContent(ctx, p.model.modelName, genai.Text(p.simplePrompt), config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}
	return &GeminiResponse{
		model:    p.model,
		response: response,
		text:     lazy.New(response.Text),
	}, nil
}

func (p *PromptBuilder) WithSchema(schema fluentllm.Schema) fluentllm.PromptBuilder {
	geminiSchema, ok := schema.(*Schema)
	if !ok {
		panic(fmt.Sprintf("expected schema of type *Schema, got %T", schema))
	}
	p.schema = geminiSchema.schema
	return p
}

func (r *GeminiResponse) Text() string {
	return r.text.Get()
}

func (r *GeminiResponse) String() string {
	return fmt.Sprintf("GeminiResponse{model=%s, text='%s'}", r.model, r.Text())
}

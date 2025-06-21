package fluentllm_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/snark87/fluentllm"
	"github.com/snark87/fluentllm/gemini"
	"github.com/snark87/fluentllm/internal/testutils"
)

func TestSimplePrompt(t *testing.T) {
	ctx := context.Background()
	geminiTestCredentials := &gemini.Credentials{
		APIKey: useGeminiAPIKey(t),
	}
	model := gemini.NewModel("gemini-1.5-flash", geminiTestCredentials)
	response := testutils.Must1(model.Prompt("What is the capital of France?").Execute(ctx))(t)
	if response == nil || response.Text() == "" {
		t.Fatalf("Expected a valid response, got `%v`", response)
	}
	if !strings.Contains(response.Text(), "Paris") {
		t.Fatalf("Expected response to containt 'Paris', got '%s'", response.Text())
	}
}

type CityInfo struct {
	Country    string `json:"country"`
	City       string `json:"city"`
	Population int    `json:"population"`
	PhoneCode  string `json:"phone_code"`
}

func TestStructuredResponse_FromGoStruct(t *testing.T) {
	ctx := context.Background()
	geminiTestCredentials := &gemini.Credentials{
		APIKey: useGeminiAPIKey(t),
	}
	model := gemini.NewModel("gemini-1.5-flash", geminiTestCredentials)
	schema := model.NewSchema()
	responseSchema := testutils.Must1(schema.FromGoValue(CityInfo{}))(t)
	response := testutils.Must1(model.
		Prompt("What is the capital of France?").
		WithSchema(responseSchema).
		Execute(ctx),
	)(t)
	output := testutils.Must1(fluentllm.AsStructuredResponse[CityInfo](response))(t)
	if output.Country != "France" {
		t.Fatalf("Expected country to be 'France', got '%s'", output.Country)
	}
	if output.City != "Paris" {
		t.Fatalf("Expected city to be 'Paris', got '%s'", output.City)
	}
	if output.Population <= 0 {
		t.Fatalf("Expected population to be greater than 0, got %d", output.Population)
	}
	if output.PhoneCode != "+33" {
		t.Fatalf("Expected phone code to be '+33', got '%s'", output.PhoneCode)
	}
}

func TestStructuredResponse_Fluent(t *testing.T) {
	ctx := context.Background()
	geminiTestCredentials := &gemini.Credentials{
		APIKey: useGeminiAPIKey(t),
	}
	model := gemini.NewModel("gemini-1.5-flash", geminiTestCredentials)
	schema := model.NewSchema()
	responseSchema := schema.Object("response",
		schema.Str("country", schema.Required()),
		schema.Str("city", schema.Required()),
		schema.Int("population", schema.Required()),
		schema.Str("phone_code", schema.Required()),
	)

	response := testutils.Must1(model.
		Prompt("What is the capital of France?").
		WithSchema(responseSchema).
		Execute(ctx),
	)(t)
	output := testutils.Must1(fluentllm.AsStructuredResponse[CityInfo](response))(t)
	if output.Country != "France" {
		t.Fatalf("Expected country to be 'France', got '%s'", output.Country)
	}
	if output.City != "Paris" {
		t.Fatalf("Expected city to be 'Paris', got '%s'", output.City)
	}
	if output.Population <= 0 {
		t.Fatalf("Expected population to be greater than 0, got %d", output.Population)
	}
	if output.PhoneCode != "+33" {
		t.Fatalf("Expected phone code to be '+33', got '%s'", output.PhoneCode)
	}
}

func useGeminiAPIKey(t *testing.T) string {
	t.Helper()
	testutils.Must0(godotenv.Load())
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY not set, skipping test")
	}

	return apiKey
}

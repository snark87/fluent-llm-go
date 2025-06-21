# FluentLLM Go

> **⚠️ Work in Progress** - This library is currently under development

A fluent Go library for interacting with Large Language Models (LLMs) with support for structured responses and schema validation.

## Features

- **🎯 Fluent API** - Clean, chainable method calls for intuitive LLM interactions
- **📋 Structured Responses** - Automatic JSON schema generation from Go structs
- **🔒 Type Safety** - Generic `AsStructuredResponse[T]` for strongly-typed outputs
- **🔌 Provider Abstraction** - Extensible design supporting multiple LLM providers
- **✨ Developer Experience** - Simple credential management and error handling

## Currently Supported Providers

- ✅ **Google Gemini**
- 🚧 More providers coming soon

## Quick Start

### Installation

```bash
go get github.com/snark87/fluentllm
```

### Basic Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/snark87/fluentllm/gemini"
)

func main() {
    ctx := context.Background()

    // Initialize Gemini model
    credentials := &gemini.Credentials{
        APIKey: "your-gemini-api-key", // pragma: allowlist secret
    }
    model := gemini.NewModel("gemini-1.5-flash", credentials)

    // Simple text prompt
    response, err := model.
        Prompt("What is the capital of France?").
        Execute(ctx)
    if err != nil {
        panic(err)
    }

    fmt.Println(response.Text()) // Output: Paris is the capital of France...
}
```

### Structured Responses

#### Using Go Structs

```go
type CityInfo struct {
    Country    string `json:"country"`
    City       string `json:"city"`
    Population int    `json:"population"`
    PhoneCode  string `json:"phone_code"`
}

func main() {
    // ... model initialization ...

    schema := model.NewSchema()
    responseSchema, _ := schema.FromGoValue(CityInfo{})

    response, _ := model.
        Prompt("What is the capital of France?").
        WithSchema(responseSchema).
        Execute(ctx)

    output, _ := fluentllm.AsStructuredResponse[CityInfo](response)
    fmt.Printf("City: %s, Country: %s\n", output.City, output.Country)
}
```

#### Using Fluent Schema Builder

```go
schema := model.NewSchema()
responseSchema := schema.Object("response",
    schema.Str("country", schema.Required()),
    schema.Str("city", schema.Required()),
    schema.Int("population", schema.Required()),
    schema.Str("phone_code", schema.Required()),
)

response, _ := model.
    Prompt("What is the capital of France?").
    WithSchema(responseSchema).
    Execute(ctx)
```

## Project Status

This project is in early development. Current focus areas:

- [x] Core fluent API design
- [x] Google Gemini integration
- [x] Structured response support
- [x] Schema generation from Go structs
- [x] Fluent schema builder
- [ ] Additional LLM provider support
- [ ] Advanced prompt features
- [ ] Comprehensive documentation
- [ ] Production-ready error handling

## Requirements

- Go 1.24.4 or later
- Valid API credentials for supported LLM providers

## Environment Setup

For testing, create a `.env` file with:

```bash
GEMINI_API_KEY=your_gemini_api_key_here
```

## Contributing

This project is in active development. Contributions, ideas, and feedback are welcome!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

# Makefile for fluentllm - A fluent interface for Large Language Models in Go
# Project: github.com/snark87/fluentllm

# Build configuration
BUILD_DIR := ./.build
COVERAGE_FILE := $(BUILD_DIR)/coverage.out
COVERAGE_HTML := $(BUILD_DIR)/coverage.html

# Go command configuration
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOIMPORTS := goimports

# Project configuration
PACKAGE := github.com/snark87/fluentllm
MODULE_NAME := $(shell go list -m)

.PHONY: all check build clean format imports test test-watch test-coverage test-coverage-func lint vet deps tidy verify pre-commit setup-pre-commit install-tools install version help

# Default target
all: check

# Main development workflow
check: deps format imports vet lint test
	@echo "✅ All checks passed!"

# Show current version
version:
	@echo "fluentllm version: $(shell grep 'const Version' fluentllm.go | cut -d'"' -f2)"

# Formatting and code quality
format:
	@echo "🔧 Formatting Go code..."
	@which $(GOFMT) > /dev/null || (echo "❌ gofmt not found in PATH" && exit 1)
	$(GOFMT) -w -s .

imports:
	@echo "📦 Organizing imports..."
	@which $(GOIMPORTS) > /dev/null || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	$(GOIMPORTS) -w .

vet:
	@echo "🔍 Running go vet..."
	$(GOCMD) vet ./...

lint:
	@echo "🧹 Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Build and clean
build:
	@echo "🏗️ Building library (validation)..."
	$(GOBUILD) ./...

clean:
	@echo "🧽 Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Testing
test:
	@echo "🧪 Running tests..."
	@which gotestsum > /dev/null || (echo "Installing gotestsum..." && go install gotest.tools/gotestsum@latest)
	gotestsum --format testname ./...

test-watch:
	@echo "👀 Running tests in watch mode..."
	@which gotestsum > /dev/null || (echo "Installing gotestsum..." && go install gotest.tools/gotestsum@latest)
	gotestsum --watch --format testname ./...

test-coverage:
	@echo "📊 Running tests with coverage..."
	@mkdir -p $(BUILD_DIR)
	$(GOTEST) -v -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated at $(COVERAGE_HTML)"

test-coverage-func:
	@echo "📈 Running tests with function coverage..."
	@mkdir -p $(BUILD_DIR)
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE)

# Dependencies management
deps:
	@echo "📥 Downloading dependencies..."
	$(GOMOD) download

tidy:
	@echo "🧹 Tidying dependencies..."
	$(GOMOD) tidy

verify:
	@echo "✅ Verifying dependencies..."
	$(GOMOD) verify

pre-commit: setup-pre-commit
	@echo "🚀 Running pre-commit hooks..."
	@which pre-commit > /dev/null || (echo "Installing pre-commit..." && pipx install pre-commit)
	pre-commit run --all-files

setup-pre-commit:
	@echo "⚙️ Setting up pre-commit..."
	@chmod +x setup-precommit.sh
	@./setup-precommit.sh

install-tools:
	@echo "🔧 Installing development tools..."
	go install gotest.tools/gotestsum@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

install: deps build
	@echo "✅ Library validated successfully!"

# Development helpers
dev-setup: install-tools setup-pre-commit
	@echo "🚀 Development environment setup complete!"

quick-test:
	@echo "⚡ Running quick tests..."
	$(GOTEST) -short ./...

benchmark:
	@echo "📊 Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

doc:
	@echo "📚 Generating documentation..."
	@which godoc > /dev/null || (echo "Installing godoc..." && go install golang.org/x/tools/cmd/godoc@latest)
	@echo "Documentation server starting at http://localhost:6060/pkg/$(PACKAGE)/"
	godoc -http=:6060

help:
	@echo "Available targets:"
	@echo "  all            - Run all checks (format, lint, test) [default]"
	@echo "  check          - Run all checks (format, lint, test)"
	@echo "  build          - Build/validate the library"
	@echo "  clean          - Clean build artifacts"
	@echo ""
	@echo "Code Quality:"
	@echo "  format         - Format Go code with gofmt"
	@echo "  imports        - Organize imports with goimports"
	@echo "  lint           - Run golangci-lint"
	@echo "  vet            - Run go vet"
	@echo ""
	@echo "Testing:"
	@echo "  test           - Run tests with gotestsum"
	@echo "  test-watch     - Run tests in watch mode"
	@echo "  test-coverage  - Run tests with HTML coverage report"
	@echo "  test-coverage-func - Run tests with function coverage summary"
	@echo "  quick-test     - Run quick tests (short mode)"
	@echo "  benchmark      - Run benchmarks"
	@echo ""
	@echo "Dependencies:"
	@echo "  deps           - Download dependencies"
	@echo "  tidy           - Tidy up dependencies"
	@echo "  verify         - Verify dependencies"
	@echo ""
	@echo "Development:"
	@echo "  dev-setup      - Complete development environment setup"
	@echo "  pre-commit     - Run pre-commit hooks on all files"
	@echo "  setup-pre-commit - Set up pre-commit framework"
	@echo "  install-tools  - Install development tools"
	@echo "  install        - Download deps and validate build"
	@echo "  doc            - Start documentation server"
	@echo "  version        - Show current version"
	@echo "  help           - Show this help message"

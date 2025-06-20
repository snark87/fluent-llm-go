#!/bin/bash
# Pre-commit setup script for Go project

set -e

echo "🚀 Setting up pre-commit framework..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Install required Go tools
echo "📦 Installing Go tools..."

# goimports for import formatting
if ! command -v goimports &> /dev/null; then
    echo "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
fi

# gocyclo for cyclomatic complexity checking
if ! command -v gocyclo &> /dev/null; then
    echo "Installing gocyclo..."
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
fi

# golangci-lint for comprehensive linting
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

# Check if pre-commit is installed
if ! command -v pre-commit &> /dev/null; then
    echo "❌ pre-commit is not installed. Please install it first:"
    echo "   pipx install pre-commit"
    echo "   OR"
    echo "   pip install pre-commit"
    exit 1
fi

# Check if detect-secrets is installed
if ! command -v detect-secrets &> /dev/null; then
    echo "❌ detect-secrets is not installed. Please install it first:"
    echo "   pipx install detect-secrets"
    exit 1
fi

# Create secrets baseline if it doesn't exist
if [ ! -f ".secrets.baseline" ]; then
    echo "🔐 Creating secrets baseline..."
    detect-secrets scan --baseline .secrets.baseline
fi

# Install pre-commit hooks
echo "🔧 Installing pre-commit hooks..."
pre-commit install
pre-commit install --hook-type commit-msg
pre-commit install --hook-type pre-push

# Migrate config if needed
echo "🔄 Migrating config..."
pre-commit migrate-config 2>/dev/null || true

# Run hooks on all files to test
echo "🧪 Testing hooks on all files..."
pre-commit run --all-files

echo "✅ Pre-commit framework setup complete!"
echo ""
echo "📋 What's been set up:"
echo "   ✓ Pre-commit hooks for code quality"
echo "   ✓ Go formatting and linting"
echo "   ✓ Security scanning for secrets"
echo "   ✓ Commit message validation"
echo "   ✓ YAML and file checks"
echo ""
echo "🎯 Next steps:"
echo "   1. Make your first commit to test the hooks"
echo "   2. Follow conventional commit format: feat/fix/docs/etc: description"
echo "   3. Hooks will run automatically on each commit"
echo "   4. Run 'pre-commit run --all-files' to test manually"

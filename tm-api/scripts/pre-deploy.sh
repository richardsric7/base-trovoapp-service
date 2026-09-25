#!/bin/bash
# Pre-deployment script to ensure swagger docs are generated and fixed
# This should be run before building/deploying to any environment

set -e  # Exit on error

echo "🔧 Generating and fixing Swagger documentation..."

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
    export PATH=$PATH:$(go env GOPATH)/bin
fi

# Generate swagger docs
echo "Generating Swagger docs..."
make swagger-clean

# Verify the structure is correct
if ! head -1 docs/swagger.yaml | grep -q "^swagger:"; then
    echo "❌ Error: swagger.yaml structure is incorrect after generation"
    echo "Attempting to fix..."
    ./scripts/fix-swagger-yaml.sh
    
    # Verify again
    if ! head -1 docs/swagger.yaml | grep -q "^swagger:"; then
        echo "❌ Error: Failed to fix swagger.yaml structure"
        exit 1
    fi
fi

echo "✅ Swagger documentation generated and verified successfully"


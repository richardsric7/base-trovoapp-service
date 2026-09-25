#!/bin/bash
# Script to generate swagger docs in Docker build
set -e

# Find swag binary
SWAG_PATH=$(find /root -name swag 2>/dev/null | head -1)
if [ -z "$SWAG_PATH" ]; then
    SWAG_PATH=$(find /go -type f -name swag 2>/dev/null | head -1)
fi

if [ -z "$SWAG_PATH" ]; then
    echo "Error: swag binary not found"
    exit 1
fi

echo "Using swag at: $SWAG_PATH"

# Generate swagger docs
rm -rf docs/
$SWAG_PATH init

# Fix YAML structure
chmod +x scripts/fix-swagger-yaml.sh
bash scripts/fix-swagger-yaml.sh

echo "Swagger docs generated and fixed successfully"


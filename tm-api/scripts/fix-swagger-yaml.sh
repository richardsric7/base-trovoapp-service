#!/bin/bash
# Fix swagger.yaml structure by moving swagger and info to the top
# This fixes a known bug in swag where it generates YAML with wrong structure

if [ ! -f "docs/swagger.yaml" ]; then
    echo "Error: docs/swagger.yaml not found"
    exit 1
fi

# Find line numbers
swagger_line=$(grep -n "^swagger:" docs/swagger.yaml | cut -d: -f1)
info_line=$(grep -n "^info:" docs/swagger.yaml | cut -d: -f1)

if [ -z "$swagger_line" ] || [ -z "$info_line" ]; then
    echo "Warning: Could not find swagger or info lines"
    exit 0
fi

# Check if structure is already correct
if [ "$swagger_line" -lt "$info_line" ]; then
    echo "swagger.yaml structure is already correct"
    exit 0
fi

echo "Fixing swagger.yaml structure..."

# Create temporary file with correct structure
# 1. swagger version
sed -n "${swagger_line}p" docs/swagger.yaml > /tmp/swagger_fix.yaml

# 2. info section
sed -n "${info_line},$((swagger_line - 1))p" docs/swagger.yaml >> /tmp/swagger_fix.yaml

# 3. Everything before info
sed -n "1,$((info_line - 1))p" docs/swagger.yaml >> /tmp/swagger_fix.yaml

# 4. Everything after swagger
sed -n "$((swagger_line + 1)),\$p" docs/swagger.yaml >> /tmp/swagger_fix.yaml

# Replace original file
mv /tmp/swagger_fix.yaml docs/swagger.yaml

echo "swagger.yaml structure fixed successfully"


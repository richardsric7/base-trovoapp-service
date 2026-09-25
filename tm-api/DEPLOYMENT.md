# Deployment Guide - Swagger Documentation Fix

This guide explains how to ensure the Swagger YAML structure fix runs in your CI/CD pipeline.

## Problem

The `swag` tool generates `swagger.yaml` with incorrect structure (swagger version at the end instead of beginning), causing parser errors on deployment.

## Solution

We've created a post-generation script (`scripts/fix-swagger-yaml.sh`) that automatically fixes the YAML structure. This script must run after `swag init` in your build/deployment pipeline.

## Integration Methods

### 1. Docker Build (Recommended)

The Dockerfile has been updated to automatically:
- Install `swag`
- Generate swagger docs
- Run the fix script

**No additional steps needed** - it's already integrated!

### 2. GitHub Actions

If using GitHub Actions, use the provided `.github/workflows/deploy-dev.yml`:

```yaml
- name: Generate Swagger docs
  run: |
    export PATH=$PATH:$(go env GOPATH)/bin
    make swagger-clean
```

The Makefile automatically runs the fix script after `swag init`.

### 3. GitLab CI

Copy `.gitlab-ci.yml.example` to `.gitlab-ci.yml` and customize:

```yaml
script:
  - make swagger-clean
  - # Verify structure
```

### 4. Jenkins / Other CI/CD

Add this to your build script:

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate and fix swagger docs
export PATH=$PATH:$(go env GOPATH)/bin
make swagger-clean

# Or manually:
# swag init
# ./scripts/fix-swagger-yaml.sh
```

### 5. Pre-deployment Script

Use `scripts/pre-deploy.sh` as a pre-deployment hook:

```bash
./scripts/pre-deploy.sh
```

This script:
- Checks if swag is installed
- Generates swagger docs
- Runs the fix script
- Verifies the structure is correct

## Verification

To verify the fix worked, check that `swagger.yaml` starts with:

```yaml
swagger: "2.0"
info:
  ...
```

You can verify in CI/CD:

```bash
if ! head -1 docs/swagger.yaml | grep -q "^swagger:"; then
  echo "Error: swagger.yaml structure is incorrect"
  exit 1
fi
```

## Manual Testing

Before deploying, test locally:

```bash
make swagger-clean
head -5 docs/swagger.yaml  # Should start with "swagger:"
```

## Troubleshooting

### Script not running in CI/CD

1. Ensure the script is executable: `chmod +x scripts/fix-swagger-yaml.sh`
2. Check that the script path is correct in your CI/CD config
3. Verify `swag` is installed before running the script

### Script fails in Docker

- Ensure `bash` is installed in the Docker image
- Check that the script has execute permissions
- Verify the working directory is correct

### Still getting parser errors

- Verify the script actually ran (check CI/CD logs)
- Manually test: `./scripts/fix-swagger-yaml.sh`
- Check that `docs/swagger.yaml` exists before the script runs


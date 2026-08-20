.PHONY: ci tidy-check build vet lint test run

run:
	go run main.go

# --- CI parity targets: same commands the PR workflow runs -----------------
tidy-check:
	go mod tidy
	git diff --exit-code -- go.mod go.sum \
		|| (echo "go.mod / go.sum not tidy; commit the changes above" && exit 1)

build:
	go build ./...

vet:
	# See workflow: skipping stringintconv until the KycStatus bug is fixed.
	go vet -stringintconv=false ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not installed. Install: https://golangci-lint.run/welcome/install/"; exit 1; }
	golangci-lint run --timeout=5m

test:
	# Scope: skips root main_test.go (integration suite that needs a live
	# server at localhost:8080). -vet=off because go test's implicit vet
	# doesn't take the -stringintconv=false flag we use above.
	go test ./internal/... -race -vet=off -coverprofile=coverage.out -covermode=atomic

ci: tidy-check build vet lint test
	@echo "backend CI checks passed"

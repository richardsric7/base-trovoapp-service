.PHONY: lint run swagger swagger-clean kill ci tidy-check build vet test

lint:
	golangci-lint run -v --fix

run:
	go run main.go

swagger:
	export PATH=$$PATH:$$(go env GOPATH)/bin && swag init
	@./scripts/fix-swagger-yaml.sh

swagger-clean:
	rm -rf docs/
	export PATH=$$PATH:$$(go env GOPATH)/bin && swag init
	@./scripts/fix-swagger-yaml.sh

kill:
	@lsof -ti:$(filter-out $@,$(MAKECMDGOALS)) | xargs kill -9 2>/dev/null || true; \
	echo "✅ Port $(filter-out $@,$(MAKECMDGOALS)) cleared."

# --- CI parity targets: same commands the PR workflow runs -----------------
tidy-check:
	go mod tidy
	git diff --exit-code -- go.mod go.sum \
		|| (echo "go.mod / go.sum not tidy; commit the changes above" && exit 1)

build:
	go build ./...

vet:
	go vet ./...

test:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic

ci: tidy-check build vet lint test
	@echo "backend CI checks passed"

.EXPORT_ALL_VARIABLES:
GOBIN = $(shell pwd)/bin

.PHONY: init
init: tools

# Единая команда для ии-агента
.PHONY: check
check: lint test
	@git diff --exit-code -- ./go.mod ./go.sum || { echo "❌ Изменения в go.mod или go.sum — остановись и попроси человека посмотреть твои изменения"; exit 1; }
	@git diff --exit-code -- '**/testdata/**' || { echo "❌ Изменения в testdata — остановись и попроси человека посмотреть твои изменения"; exit 1; }
	@git ls-files --others --exclude-standard -- '**/testdata/**' | grep -q . && { echo "❌ Новые файлы в testdata — остановись и попроси человека посмотреть твои изменения"; exit 1; } || true

.PHONY: deps
deps:
	@go mod tidy

.PHONY: clean-mocks
clean-mocks:
	@./scripts/clean_mocks.sh

.PHONY: mocks
mocks: tools-mocks clean-mocks
	@export PATH="$(shell pwd)/bin:$(PATH)"; mockery --log-level=warn --config=.mockery.yaml

.PHONY: audit
audit: tools-govulncheck
	@export PATH="$(shell pwd)/bin:$(PATH)"; govulncheck ./...

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@go test -race -failfast -count=1 ./...

.PHONY: integration-sqlite
integration-sqlite:
	@go test -tags "integration sqlite" ./internal/sqlite/... -v -p 1 -count=1 -run $(or $(RUN),TestIntegration)

.PHONY: run
run: tools-goose
	@./bin/goose up
	@go run -race ./cmd/capuchin/main.go

.PHONY: tools
tools: tools-goose tools-govulncheck tools-mocks

.PHONY: tools-goose
tools-goose:
	@go install github.com/pressly/goose/v3/cmd/goose@v3.28.0

.PHONY: tools-govulncheck
tools-govulncheck:
	@go install golang.org/x/vuln/cmd/govulncheck@v1.7.0

.PHONY: tools-mocks
tools-mocks:
	@go install github.com/vektra/mockery/v2@v2.53.7

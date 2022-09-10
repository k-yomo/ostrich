.DEFAULT_GOAL := help

.PHONY: test
test: ## Run test
	go install gotest.tools/gotestsum@latest
	gotestsum -- -race -coverprofile=coverage.out ./...

.PHONY: cover
cover: test ## Run test with showing coverage
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: fmt
fmt: ## Format code
	goimports -w .

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

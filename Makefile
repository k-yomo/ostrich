.DEFAULT_GOAL := help

.PHONY: test
test: ## Run test
	go test -race -coverprofile=coverage.out ./...

.PHONY: cover
cover: test ## Run test with showing coverage
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

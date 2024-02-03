.DEFAULT_GOAL := help

.PHONY: test
test: ## Run all tests.
	@go test -race -count=1 -coverprofile=cover.out ./...
	@go tool cover -func=cover.out | grep ^total | tr -s '\t'

.PHONY: test-short
test-short: ## Run only unit tests, tests without I/O dependencies.
	@go test -short ./...

.PHONY: test-env-up
test-env-up: ## Run test environment.
	@docker compose -f docker-compose.test.yml up -d

.PHONY: test-env-down
test-env-down: ## Down and cleanup test environment.
	@docker compose -f docker-compose.test.yml down -v

.PHONY: cover
cover: ## Show coverage in browser
	@go tool cover -html=cover.out

.PHONY: help
help: ## Show help for each of the Makefile targets.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

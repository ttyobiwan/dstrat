SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help section
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+.*:.*?##/ {printf "\033[36m%-38s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: temporal
temporal: ## Start temporal server
	temporal server start-dev --ui-port 9090 --db-filename temporal.sqlite

.PHONY: server
server: ## Start echo server using air
	air

.PHONY: worker
worker: ## Start temporal worker
	go run cmd/worker/* -q $(q)

.PHONY: test
test: ## Run go tests (without caching)
	go test ./... -count=1

# Thin wrapper around mage (https://magefile.org/); see `make help` or magefile.go.
# mage is run from vendor/, so it does not need to be installed.
MAGE := go run github.com/magefile/mage

.DEFAULT_GOAL := build

.PHONY: build
build: ## Build the countbeat binary
	$(MAGE) build

.PHONY: test
test: ## Run unit tests
	go test -race ./...

.PHONY: update
update: ## Regenerate fields.yml, include/fields.go, configs and field docs
	$(MAGE) update

.PHONY: check
check: ## Format code, regenerate files and fail on uncommitted changes
	$(MAGE) check

.PHONY: fmt
fmt: ## Format source code
	$(MAGE) fmt

.PHONY: clean
clean: ## Remove build artifacts
	$(MAGE) clean

.PHONY: package
package: ## Build distribution packages (requires Docker)
	$(MAGE) package

.PHONY: vendor
vendor: ## Tidy go.mod and refresh vendor/
	go mod tidy
	go mod vendor

.PHONY: help
help: ## Show this help
	@grep -E '^[a-z]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-10s %s\n", $$1, $$2}'

# VARIABLES
BINARY_NAME="stock"
BINARY_PATH = "build"

.PHONY: build

default: usage

clean: ## Trash binary files
	@echo "--> cleaning..."
	@rm -rf $(BINARY_PATH) 2> /dev/null
	@echo "Clean OK"

test: ## Run all tests
	@echo "--> testing..."
	@sh script/test.sh

coverage: ## Run all tests with coverage
	@echo "--> running tests with coverage"
	@sh script/coverage.sh

coverage-html: ## Run all tests with coverage opening HTML report in browser
	@echo "--> running tests with coverage, HTML report will open in browser"
	@sh script/coverage.sh --html

lint: ## Run golint
	@echo "--> running golint..."
	@sh script/lint.sh

run: ## Run your application
	@echo $(GOPATH)
	@sh script/dev.sh

migrate: ## Run your application
	@echo $(GOPATH)
	@sh script/migrate.sh

build: ## build the application
	@echo "--> building..."
	@sh script/build.sh
	@echo "Build OK"

release: ## build the application and create the artifacts
	@echo "--> building..."
	@sh script/build.sh -r
	@echo "Release OK"

usage: ## List available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

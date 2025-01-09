# Project variables
BINARY_NAME=pdfmod
SOURCE_DIR=./cmd/pdfmod/main.go
BUILD_DIR=./bin
MOCKS_DIR=./mocks

# Colors for output
BLUE=\033[0;34m
GREEN=\033[0;32m
RESET=\033[0m

# Default goal
.DEFAULT_GOAL := help

# Build the binary
build: ## Build the Go binary
	@echo "$(BLUE)Building the project...$(RESET)"
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SOURCE_DIR)
	@echo "$(GREEN)Build complete! Binary created at $(BUILD_DIR)/$(BINARY_NAME)$(RESET)"

# Run the binary
run: build ## Build and run the Go project
	@echo "$(BLUE)Running the project...$(RESET)"
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run locally
run-local: ## Run the project locally
	@echo "$(BLUE)Running the project locally...$(RESET)"
	@go run $(SOURCE_DIR)

# Clean the build
clean: ## Remove build files
	@echo "$(BLUE)Cleaning up...$(RESET)"
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@echo "$(GREEN)Cleanup complete!$(RESET)"

# Generate mocks
gen-mocks: ## Generate mocks for interfaces using mockgen
	@echo "$(BLUE)Generating mocks...$(RESET)"
	@mockgen -source=internal/file/interfaces.go -destination=$(MOCKS_DIR)/mock_file_handler.go -package=mocks
	@mockgen -source=internal/pdf/interfaces.go -destination=$(MOCKS_DIR)/mock_pdf_metadata_handler.go -package=mocks
	@mockgen -destination=$(MOCKS_DIR)/mock_file_info.go -package=mocks io/fs FileInfo
	@mockgen -destination=$(MOCKS_DIR)/mock_pdf_manager.go -package=mocks github.com/sidshirsat/pdfmod/internal/manager PDFManagerInterface
	@mockgen -destination=$(MOCKS_DIR)/mock_prompter.go -package=mocks github.com/sidshirsat/pdfmod/internal/manager Prompter
	@echo "$(GREEN)Mocks generated successfully!$(RESET)"

.PHONY: test
# Run tests
test: ## Run tests with and without race detector, generate coverage and JUnit report
	@echo "Running tests without race detector..."
	@CGO_ENABLED=0 gotestsum --format testname --junitfile junit-tests.xml -- -cover -coverprofile=coverage.out ./...
	@echo "Running tests with race detector..."
	@CGO_ENABLED=1 gotestsum --format testname -- -race ./...

# Lint the code
lint: ## Run golangci-lint on the code
	@echo "$(BLUE)Running lint checks...$(RESET)"
	@golangci-lint run
	@echo "$(GREEN)Lint checks complete!$(RESET)"

# Format the code
fmt: ## Format the code using gofmt
	@echo "$(BLUE)Formatting code...$(RESET)"
	@gofmt -s -w .
	@echo "$(GREEN)Formatting complete!$(RESET)"

# Run all checks before committing
check: fmt lint test ## Run formatting, linting, and tests

# Display help for each target
help: ## Display this help message
	@echo "$(BLUE)Available commands:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-15s$(RESET) %s\n", $$1, $$2}'


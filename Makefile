# Binary name
BINARY_NAME=git-clone-tool
VERSION?=1.0.0
BUILD_DIR=build
CMD_DIR=cmd/clone-git-repo

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BUILD_DIR)

# Build information
BUILD_TIME=$(shell date +%FT%T%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Use linker flags to provide version/build information
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.GitBranch=$(GIT_BRANCH)"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all build clean test help

all: clean build ## Build the binary

build: ## Build the binary
	@echo "Building..."
	mkdir -p $(BUILD_DIR)
	cd $(CMD_DIR) && go build $(LDFLAGS) -o ../../$(BUILD_DIR)/$(BINARY_NAME)
	@echo "Binary built at $(BUILD_DIR)/$(BINARY_NAME)"

clean: ## Remove build related files
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean
	@echo "Cleaned!"

test: ## Run tests
	@echo "Running tests..."
	go test ./... -v

lint: ## Run linting
	@echo "Running linting..."
	go vet ./...
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Please install it to run linting."; \
	fi

fmt: ## Run go fmt
	@echo "Running go fmt..."
	go fmt ./...

run: build ## Run the binary
	./$(BUILD_DIR)/$(BINARY_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Cross compilation targets
.PHONY: build-linux build-windows build-darwin build-all

build-linux: ## Build for Linux
	@echo "Building for Linux..."
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 cd $(CMD_DIR) && go build $(LDFLAGS) -o ../../$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64

build-windows: ## Build for Windows
	@echo "Building for Windows..."
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 cd $(CMD_DIR) && go build $(LDFLAGS) -o ../../$(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe

build-darwin: ## Build for macOS
	@echo "Building for macOS..."
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 cd $(CMD_DIR) && go build $(LDFLAGS) -o ../../$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64

build-all: build-linux build-windows build-darwin ## Build for all platforms

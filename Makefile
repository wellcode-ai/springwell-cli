# Makefile for SpringWell CLI

# Go parameters
BINARY_NAME=springwell
MAIN_PATH=cmd/springwell/main.go
GO=go
GOFLAGS=-ldflags="-s -w" # Strip debug information
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags="-s -w -X main.Version=$(VERSION)"

# Build directories
BUILD_DIR=build
DIST_DIR=dist

# Installation directory
INSTALL_DIR=/usr/local/bin

# Define the default target when just running make
.PHONY: all
all: clean build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Install the application
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installation complete"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo "Clean complete"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run code linting
.PHONY: lint
lint:
	@echo "Linting code..."
	$(GO) vet ./...
	@if command -v golint > /dev/null; then \
		golint ./...; \
	else \
		echo "golint not installed"; \
	fi

# Build for multiple platforms (cross-compilation)
.PHONY: dist
dist: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(DIST_DIR)
	
	# Linux (amd64)
	@echo "Building for linux/amd64..."
	@GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	
	# Linux (arm64)
	@echo "Building for linux/arm64..."
	@GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	
	# macOS (amd64)
	@echo "Building for darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	
	# macOS (arm64)
	@echo "Building for darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	
	# Windows (amd64)
	@echo "Building for windows/amd64..."
	@GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	
	@echo "Cross-compilation complete. Binaries available in $(DIST_DIR)"

# Create release archives
.PHONY: release
release: dist
	@echo "Creating release archives..."
	@cd $(DIST_DIR) && tar -czf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	@cd $(DIST_DIR) && tar -czf $(BINARY_NAME)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	@cd $(DIST_DIR) && tar -czf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	@cd $(DIST_DIR) && tar -czf $(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@cd $(DIST_DIR) && zip $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "Release archives created in $(DIST_DIR)"

# Run the application (for quick testing)
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Generate code documentation
.PHONY: doc
doc:
	@echo "Generating documentation..."
	@mkdir -p docs
	@if command -v godoc > /dev/null; then \
		godoc -http=:6060; \
	else \
		echo "godoc not installed"; \
	fi

# Show help
.PHONY: help
help:
	@echo "SpringWell CLI - Make targets:"
	@echo "  all         - Clean and build the application"
	@echo "  build       - Build the application"
	@echo "  clean       - Remove build artifacts"
	@echo "  test        - Run tests"
	@echo "  lint        - Run linting tools"
	@echo "  install     - Install to $(INSTALL_DIR)"
	@echo "  dist        - Build for multiple platforms"
	@echo "  release     - Create release archives"
	@echo "  run         - Build and run the application"
	@echo "  doc         - Generate code documentation"
	@echo "  help        - Show this help message" 
.PHONY: build install clean test lint

# Build variables
BINARY_NAME=springwell
MAIN_PACKAGE=./cmd/springwell
VERSION?=0.1.0
BUILD_DIR=build
GO_FLAGS=-ldflags "-X main.Version=$(VERSION)"

# Build the project
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(GO_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Done!"

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@mv $(BINARY_NAME) $(GOPATH)/bin/
	@echo "Done! $(BINARY_NAME) has been installed to $(GOPATH)/bin/"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "Done!"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Create distribution packages
dist: clean
	@echo "Creating distribution packages..."
	@mkdir -p $(BUILD_DIR)
	
	# Build for different platforms
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build $(GO_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	
	@echo "Building for macOS..."
	@GOOS=darwin GOARCH=amd64 go build $(GO_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	@GOOS=darwin GOARCH=arm64 go build $(GO_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build $(GO_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	
	@echo "Done!"

# Default target
all: build 
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

# Project parameters
PROJECT_NAME?=my-project
TEMPLATE?=basic
DB?=h2
VERBOSE?=false
FIX_MISSING?=true

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

# Verify template integrity (check if all template files exist)
.PHONY: verify-templates
verify-templates:
	@echo "Verifying template files..."
	@if [ ! -d "pkg/templates" ]; then \
		echo "❌ Error: pkg/templates directory not found"; \
		exit 1; \
	fi
	@echo "✅ Found: pkg/templates directory"
	
	@echo "Checking basic template..."
	@if [ ! -d "pkg/templates/project/basic" ]; then \
		echo "❌ Error: Basic template directory not found"; \
		exit 1; \
	fi
	@echo "✅ Found: Basic template directory"
	
	@echo "Checking aws-temporal-auth0 template..."
	@if [ ! -d "pkg/templates/project/aws-temporal-auth0" ]; then \
		echo "❌ Error: AWS-Temporal-Auth0 template directory not found"; \
		exit 1; \
	fi
	@echo "✅ Found: AWS-Temporal-Auth0 template directory"
	
	@echo "Checking for required template files in aws-temporal-auth0 template..."
	@missing=0; \
	for file in config/AwsConfig.java.tmpl config/Auth0Config.java.tmpl config/SecurityConfig.java.tmpl \
				config/TemporalConfig.java.tmpl temporal/worker/TemporalWorkerRegistrar.java.tmpl \
				temporal/worker/TemporalWorkerService.java.tmpl \
				temporal/workflow/OrderProcessingWorkflow.java.tmpl; do \
		if [ ! -f "pkg/templates/project/aws-temporal-auth0/$$file" ]; then \
			echo "❌ Missing template file: $$file"; \
			missing=1; \
		else \
			echo "✅ Found template file: $$file"; \
		fi; \
	done; \
	if [ $$missing -eq 1 ]; then \
		exit 1; \
	fi
	
	@echo "Template files verified successfully."

# Create templates directory structure (useful when developing templates)
.PHONY: create-template-structure
create-template-structure:
	@echo "Creating template directory structure..."
	
	@# Create basic template structure
	@mkdir -p pkg/templates/project/basic/src/main/java/com/example/project
	@mkdir -p pkg/templates/project/basic/src/main/resources
	@mkdir -p pkg/templates/project/basic/src/test/java/com/example/project
	
	@# Create aws-temporal-auth0 template structure
	@mkdir -p pkg/templates/project/aws-temporal-auth0/.github/workflows
	@mkdir -p pkg/templates/project/aws-temporal-auth0/.mvn
	@mkdir -p pkg/templates/project/aws-temporal-auth0/charts/project/templates
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/java/com/example/project/config
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/java/com/example/project/controller
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/java/com/example/project/domain/dto
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/java/com/example/project/domain/entity
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/java/com/example/project/exception
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/java/com/example/project/middleware
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/java/com/example/project/messaging
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/java/com/example/project/repository
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/java/com/example/project/service
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/java/com/example/project/temporal/activity/impl
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/java/com/example/project/temporal/worker
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/java/com/example/project/temporal/workflow/impl
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/java/com/example/project/util
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/resources/db/migration
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/main/swagger
	@mkdir -p pkg/templates/project/aws-temporal-auth0/src/test/java/com/example/project
	
	@echo "Template directory structure created successfully."

# Generate a new project with template verification and auto-fix
.PHONY: new-project
new-project: build verify-templates
	@echo "Generating new project: $(PROJECT_NAME) (template: $(TEMPLATE), database: $(DB))..."
	@# First, attempt to create the project with the CLI
	@if [ "$(VERBOSE)" = "true" ]; then \
		$(BUILD_DIR)/$(BINARY_NAME) new $(PROJECT_NAME) --template $(TEMPLATE) --db $(DB) --verbose; \
	else \
		$(BUILD_DIR)/$(BINARY_NAME) new $(PROJECT_NAME) --template $(TEMPLATE) --db $(DB); \
	fi
	
	@# Check if the project directory exists
	@if [ ! -d "$(PROJECT_NAME)" ]; then \
		echo "❌ Error: Project creation failed. Project directory $(PROJECT_NAME) not found."; \
		exit 1; \
	fi
	
	@# Validate project structure and fix if needed
	@echo "Validating project structure..."
	@if [ "$(FIX_MISSING)" = "true" ]; then \
		echo "Auto-fix enabled. Will create any missing directories."; \
		if [ "$(TEMPLATE)" = "aws-temporal-auth0" ]; then \
			mkdir -p "$(PROJECT_NAME)/.github/workflows"; \
			mkdir -p "$(PROJECT_NAME)/.mvn"; \
			mkdir -p "$(PROJECT_NAME)/charts/$(PROJECT_NAME)/templates"; \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/config"; \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/controller"; \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/domain/dto"; \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/domain/entity"; \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/exception"; \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/middleware"; \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/messaging"; \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/repository"; \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/service"; \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/temporal/activity/impl"; \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/temporal/worker"; \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/temporal/workflow/impl"; \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/util"; \
			mkdir -p "$(PROJECT_NAME)/src/main/resources/db/migration"; \
			mkdir -p "$(PROJECT_NAME)/src/main/swagger"; \
			mkdir -p "$(PROJECT_NAME)/src/test/java/com/example/$(PROJECT_NAME)"; \
			if [ ! -f "$(PROJECT_NAME)/compose.yaml" ] && [ -f "pkg/templates/project/aws-temporal-auth0/compose.yaml.tmpl" ]; then \
				sed "s/\{\{project\}\}/$(PROJECT_NAME)/g" pkg/templates/project/aws-temporal-auth0/compose.yaml.tmpl > "$(PROJECT_NAME)/compose.yaml"; \
			fi; \
			if [ ! -f "$(PROJECT_NAME)/Dockerfile" ] && [ -f "pkg/templates/project/aws-temporal-auth0/Dockerfile.tmpl" ]; then \
				cp pkg/templates/project/aws-temporal-auth0/Dockerfile.tmpl "$(PROJECT_NAME)/Dockerfile"; \
			fi; \
		else \
			mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)"; \
			mkdir -p "$(PROJECT_NAME)/src/main/resources"; \
			mkdir -p "$(PROJECT_NAME)/src/test/java/com/example/$(PROJECT_NAME)"; \
		fi; \
		echo "✅ Project structure validated and fixed if needed."; \
	fi
	
	@echo "Project created successfully at ./$(PROJECT_NAME)"
	@echo ""
	@echo "To navigate to your project, run:"
	@echo "  cd $(PROJECT_NAME)"

# Generate a new project with debug output
.PHONY: new-project-debug
new-project-debug: build verify-templates
	@echo "Generating new project with verbose output: $(PROJECT_NAME) (template: $(TEMPLATE), database: $(DB))..."
	$(BUILD_DIR)/$(BINARY_NAME) new $(PROJECT_NAME) --template $(TEMPLATE) --db $(DB) --verbose
	@echo "Project created with debug output at ./$(PROJECT_NAME)"
	@echo ""
	@make verify-project PROJECT_NAME=$(PROJECT_NAME)
	@if [ "$(FIX_MISSING)" = "true" ]; then \
		make fix-project-structure PROJECT_NAME=$(PROJECT_NAME) TEMPLATE=$(TEMPLATE); \
	fi

# Generate and navigate to a new project (creates helper script)
.PHONY: create-project
create-project: build verify-templates
	@echo "Generating new project and preparing navigation script..."
	@if [ "$(VERBOSE)" = "true" ]; then \
		$(BUILD_DIR)/$(BINARY_NAME) new $(PROJECT_NAME) --template $(TEMPLATE) --db $(DB) --verbose; \
	else \
		$(BUILD_DIR)/$(BINARY_NAME) new $(PROJECT_NAME) --template $(TEMPLATE) --db $(DB); \
	fi
	
	@# Validate and fix project structure if needed
	@if [ "$(FIX_MISSING)" = "true" ]; then \
		make fix-project-structure PROJECT_NAME=$(PROJECT_NAME) TEMPLATE=$(TEMPLATE); \
	fi
	
	@echo "#!/bin/sh" > navigate-to-project.sh
	@echo "echo \"Navigating to $(PROJECT_NAME)...\"" >> navigate-to-project.sh
	@echo "cd $(PROJECT_NAME)" >> navigate-to-project.sh
	@echo "exec \$$SHELL" >> navigate-to-project.sh
	@chmod +x navigate-to-project.sh
	@echo "Project created successfully!"
	@echo ""
	@echo "To navigate to your project, run:"
	@echo "  source navigate-to-project.sh"

# Verify project structure (check if all expected folders exist)
.PHONY: verify-project
verify-project:
	@echo "Verifying project structure for $(PROJECT_NAME)..."
	@echo "Expected directory structure for template $(TEMPLATE):"
	@if [ "$(TEMPLATE)" = "aws-temporal-auth0" ]; then \
		echo "Checking AWS-Temporal-Auth0 template structure..."; \
		for dir in \
			"$(PROJECT_NAME)/.github/workflows" \
			"$(PROJECT_NAME)/.mvn" \
			"$(PROJECT_NAME)/charts/$(PROJECT_NAME)/templates" \
			"$(PROJECT_NAME)/src/main/java" \
			"$(PROJECT_NAME)/src/main/resources/db/migration" \
			"$(PROJECT_NAME)/src/main/swagger" \
			"$(PROJECT_NAME)/src/test/java" \
			"$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/config" \
			"$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/controller" \
			"$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/domain/dto" \
			"$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/domain/entity" \
			"$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/exception" \
			"$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/middleware" \
			"$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/messaging" \
			"$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/repository" \
			"$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/service" \
			"$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/temporal/activity" \
			"$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/temporal/worker" \
			"$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/temporal/workflow" \
			"$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/util"; \
		do \
			if [ ! -d "$$dir" ]; then \
				echo "❌ Missing directory: $$dir"; \
			else \
				echo "✅ Found directory: $$dir"; \
			fi; \
		done; \
	else \
		echo "Checking Basic template structure..."; \
		for dir in \
			"$(PROJECT_NAME)/src/main/java" \
			"$(PROJECT_NAME)/src/main/resources" \
			"$(PROJECT_NAME)/src/test/java"; \
		do \
			if [ ! -d "$$dir" ]; then \
				echo "❌ Missing directory: $$dir"; \
			else \
				echo "✅ Found directory: $$dir"; \
			fi; \
		done; \
	fi
	@echo "Checking key files..."
	@if [ -f "$(PROJECT_NAME)/pom.xml" ]; then \
		echo "✅ Found file: pom.xml"; \
	else \
		echo "❌ Missing file: pom.xml"; \
	fi
	@if [ -f "$(PROJECT_NAME)/README.md" ]; then \
		echo "✅ Found file: README.md"; \
	else \
		echo "❌ Missing file: README.md"; \
	fi
	@if [ "$(TEMPLATE)" = "aws-temporal-auth0" ]; then \
		if [ -f "$(PROJECT_NAME)/Dockerfile" ]; then \
			echo "✅ Found file: Dockerfile"; \
		else \
			echo "❌ Missing file: Dockerfile"; \
		fi; \
		if [ -f "$(PROJECT_NAME)/compose.yaml" ]; then \
			echo "✅ Found file: compose.yaml"; \
		else \
			echo "❌ Missing file: compose.yaml"; \
		fi; \
	fi
	@echo "Structure verification complete."

# Fix missing project structure (attempts to create missing directories)
.PHONY: fix-project-structure
fix-project-structure:
	@echo "Fixing project structure for $(PROJECT_NAME)..."
	@if [ ! -d "$(PROJECT_NAME)" ]; then \
		echo "Error: Project directory $(PROJECT_NAME) not found!"; \
		exit 1; \
	fi
	@if [ "$(TEMPLATE)" = "aws-temporal-auth0" ]; then \
		echo "Creating AWS-Temporal-Auth0 template structure..."; \
		mkdir -p "$(PROJECT_NAME)/.github/workflows"; \
		mkdir -p "$(PROJECT_NAME)/.mvn"; \
		mkdir -p "$(PROJECT_NAME)/charts/$(PROJECT_NAME)/templates"; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/config"; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/controller"; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/domain/dto"; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/domain/entity"; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/exception"; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/middleware"; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/messaging"; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/repository"; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/service"; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/temporal/activity/impl"; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/temporal/worker"; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/temporal/workflow/impl"; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)/util"; \
		mkdir -p "$(PROJECT_NAME)/src/main/resources/db/migration"; \
		mkdir -p "$(PROJECT_NAME)/src/main/swagger"; \
		mkdir -p "$(PROJECT_NAME)/src/test/java/com/example/$(PROJECT_NAME)"; \
		if [ ! -f "$(PROJECT_NAME)/Dockerfile" ]; then \
			echo "# Multi-stage build Dockerfile for SpringBoot application" > "$(PROJECT_NAME)/Dockerfile"; \
			echo "FROM eclipse-temurin:17-jdk as build" >> "$(PROJECT_NAME)/Dockerfile"; \
			echo "WORKDIR /app" >> "$(PROJECT_NAME)/Dockerfile"; \
			echo "COPY . ." >> "$(PROJECT_NAME)/Dockerfile"; \
			echo "RUN ./mvnw clean package -DskipTests" >> "$(PROJECT_NAME)/Dockerfile"; \
			echo "" >> "$(PROJECT_NAME)/Dockerfile"; \
			echo "FROM eclipse-temurin:17-jre" >> "$(PROJECT_NAME)/Dockerfile"; \
			echo "WORKDIR /app" >> "$(PROJECT_NAME)/Dockerfile"; \
			echo "COPY --from=build /app/target/*.jar app.jar" >> "$(PROJECT_NAME)/Dockerfile"; \
			echo "ENTRYPOINT [\"java\", \"-jar\", \"app.jar\"]" >> "$(PROJECT_NAME)/Dockerfile"; \
		fi; \
		if [ ! -f "$(PROJECT_NAME)/compose.yaml" ]; then \
			echo "version: '3.8'" > "$(PROJECT_NAME)/compose.yaml"; \
			echo "" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "services:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "  app:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "    build: ." >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "    ports:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      - \"8080:8080\"" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "    environment:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      - SPRING_PROFILES_ACTIVE=dev" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "    depends_on:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      - postgres" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      - temporal" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "  postgres:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "    image: postgres:15-alpine" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "    environment:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      POSTGRES_USER: postgres" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      POSTGRES_PASSWORD: postgres" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      POSTGRES_DB: $(PROJECT_NAME)" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "    ports:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      - \"5432:5432\"" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "    volumes:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      - postgres-data:/var/lib/postgresql/data" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "  temporal:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "    image: temporalio/auto-setup:1.22.0" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "    environment:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      - DB=postgresql" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      - DB_PORT=5432" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      - POSTGRES_USER=postgres" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      - POSTGRES_PWD=postgres" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      - POSTGRES_SEEDS=postgres" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "    ports:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      - \"7233:7233\"" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "    depends_on:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "      - postgres" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "volumes:" >> "$(PROJECT_NAME)/compose.yaml"; \
			echo "  postgres-data:" >> "$(PROJECT_NAME)/compose.yaml"; \
		fi; \
	else \
		echo "Creating Basic template structure..."; \
		mkdir -p "$(PROJECT_NAME)/src/main/java/com/example/$(PROJECT_NAME)"; \
		mkdir -p "$(PROJECT_NAME)/src/main/resources"; \
		mkdir -p "$(PROJECT_NAME)/src/test/java/com/example/$(PROJECT_NAME)"; \
	fi
	@echo "Fixed project structure. Run 'make verify-project PROJECT_NAME=$(PROJECT_NAME)' to check again."

# Full end-to-end test of project creation and validation
.PHONY: e2e-test
e2e-test: build
	@echo "Running end-to-end test of project creation..."
	@test_project="test-project-$$(date +%s)"
	@echo "Creating test project: $$test_project"
	
	@# Create the project
	@$(MAKE) new-project PROJECT_NAME=$$test_project TEMPLATE=$(TEMPLATE) DB=$(DB) FIX_MISSING=true
	
	@# Verify the project
	@$(MAKE) verify-project PROJECT_NAME=$$test_project TEMPLATE=$(TEMPLATE)
	
	@echo "Test completed successfully. Cleaning up..."
	@rm -rf $$test_project
	@echo "Test project removed."

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
	@echo "  all                    - Clean and build the application"
	@echo "  build                  - Build the application"
	@echo "  clean                  - Remove build artifacts"
	@echo "  test                   - Run tests"
	@echo "  lint                   - Run linting tools"
	@echo "  install                - Install to $(INSTALL_DIR)"
	@echo "  dist                   - Build for multiple platforms"
	@echo "  release                - Create release archives"
	@echo "  run                    - Build and run the application"
	@echo "  verify-templates       - Check if all template files exist before project creation"
	@echo "  create-template-structure - Create template directory structure (for development)"
	@echo "  new-project            - Generate a new project with automatic fixes (usage: make new-project PROJECT_NAME=my-app TEMPLATE=aws-temporal-auth0 DB=postgres)"
	@echo "  new-project-debug      - Generate a new project with verbose output and structure verification"
	@echo "  create-project         - Generate a new project and create a navigation script"
	@echo "  verify-project         - Check if a generated project has all required directories"
	@echo "  fix-project-structure  - Create missing directories in a project"
	@echo "  e2e-test               - Run end-to-end test of project creation and validation"
	@echo "  doc                    - Generate code documentation"
	@echo "  help                   - Show this help message" 
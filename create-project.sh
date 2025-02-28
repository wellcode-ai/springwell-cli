#!/bin/bash

# create-project.sh - A helper script to create a SpringWell project and navigate to it

# Default values
PROJECT_NAME="my-project"
TEMPLATE="basic"
DB="h2"
CLI_PATH="./build/springwell"
VERBOSE=false
FIX_STRUCTURE=true  # Auto-fix is enabled by default
VERIFY_TEMPLATES=true # Verify templates before creating the project

# Display help
function show_help {
    echo "SpringWell Project Creator"
    echo ""
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -n, --name NAME       Project name (default: my-project)"
    echo "  -t, --template NAME   Project template (default: basic, options: basic, aws-temporal-auth0)"
    echo "  -d, --db TYPE         Database type (default: h2, options: h2, postgres, mysql)"
    echo "  -p, --path PATH       Path to the springwell CLI (default: ./build/springwell)"
    echo "  -v, --verbose         Enable verbose output"
    echo "  -f, --fix             Fix directory structure if some folders are missing (default: enabled)"
    echo "  --no-fix              Disable auto-fixing of directory structure"
    echo "  --no-verify           Disable template verification"
    echo "  -h, --help            Show this help message"
    echo ""
    echo "Example:"
    echo "  $0 --name my-awesome-app --template aws-temporal-auth0 --db postgres"
    echo ""
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -n|--name)
            PROJECT_NAME="$2"
            shift # past argument
            shift # past value
            ;;
        -t|--template)
            TEMPLATE="$2"
            shift # past argument
            shift # past value
            ;;
        -d|--db)
            DB="$2"
            shift # past argument
            shift # past value
            ;;
        -p|--path)
            CLI_PATH="$2"
            shift # past argument
            shift # past value
            ;;
        -v|--verbose)
            VERBOSE=true
            shift # past argument
            ;;
        -f|--fix)
            FIX_STRUCTURE=true
            shift # past argument
            ;;
        --no-fix)
            FIX_STRUCTURE=false
            shift # past argument
            ;;
        --no-verify)
            VERIFY_TEMPLATES=false
            shift # past argument
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Function to check if the CLI exists and is executable
check_cli() {
    if [ ! -f "$CLI_PATH" ]; then
        echo "Error: SpringWell CLI not found at $CLI_PATH"
        echo "Please build the CLI first with 'make build' or specify the correct path with --path"
        return 1
    fi

    if [ ! -x "$CLI_PATH" ]; then
        echo "Warning: CLI file is not executable. Attempting to make it executable..."
        chmod +x "$CLI_PATH"
        if [ $? -ne 0 ]; then
            echo "Error: Could not make CLI executable. Please check permissions."
            return 1
        fi
    fi

    return 0
}

# Function to verify template files existence
verify_templates() {
    local template="$1"
    
    echo "Verifying template files..."
    
    # Check templates directory existence
    if [ ! -d "$(dirname "$CLI_PATH")/pkg/templates" ]; then
        echo "❌ Error: templates directory not found"
        return 1
    fi
    
    echo "✅ Found templates directory"
    
    # Check template directory existence
    local template_dir="$(dirname "$CLI_PATH")/pkg/templates/project/$template"
    if [ ! -d "$template_dir" ]; then
        echo "❌ Error: Template directory for '$template' not found"
        echo "Available templates:"
        ls "$(dirname "$CLI_PATH")/pkg/templates/project/" 2>/dev/null || echo "  None available"
        return 1
    fi
    
    echo "✅ Found template directory for '$template'"
    
    # For aws-temporal-auth0 template, check key files
    if [ "$template" = "aws-temporal-auth0" ]; then
        local missing=0
        local required_files=(
            "config/AwsConfig.java.tmpl"
            "config/Auth0Config.java.tmpl"
            "config/SecurityConfig.java.tmpl"
            "config/TemporalConfig.java.tmpl"
            "temporal/worker/TemporalWorkerRegistrar.java.tmpl"
            "temporal/worker/TemporalWorkerService.java.tmpl"
            "temporal/workflow/OrderProcessingWorkflow.java.tmpl"
        )
        
        for file in "${required_files[@]}"; do
            if [ ! -f "$template_dir/$file" ]; then
                echo "❌ Missing required template file: $file"
                missing=1
            else
                echo "✅ Found template file: $file"
            fi
        done
        
        if [ $missing -eq 1 ]; then
            echo "Some required template files are missing!"
            return 1
        fi
    fi
    
    echo "✅ Template verification completed successfully."
    return 0
}

# Verify project structure
verify_project_structure() {
    local project="$1"
    local template="$2"
    
    echo "Verifying project structure..."
    
    # Check base project directory
    if [ ! -d "$project" ]; then
        echo "❌ Project directory not created: $project"
        return 1
    fi
    
    if [ "$template" = "aws-temporal-auth0" ]; then
        # Check key directories for AWS template
        local dirs=(
            ".github/workflows"
            ".mvn"
            "charts/$project/templates"
            "src/main/java"
            "src/main/resources/db/migration"
            "src/main/swagger"
            "src/test/java"
            "src/main/java/com/example/$project/config"
            "src/main/java/com/example/$project/controller"
            "src/main/java/com/example/$project/domain/dto"
            "src/main/java/com/example/$project/domain/entity"
            "src/main/java/com/example/$project/exception"
            "src/main/java/com/example/$project/middleware"
            "src/main/java/com/example/$project/messaging"
            "src/main/java/com/example/$project/repository"
            "src/main/java/com/example/$project/service"
            "src/main/java/com/example/$project/temporal/activity"
            "src/main/java/com/example/$project/temporal/worker"
            "src/main/java/com/example/$project/temporal/workflow"
            "src/main/java/com/example/$project/util"
        )
        
        local missing=false
        for dir in "${dirs[@]}"; do
            if [ ! -d "$project/$dir" ]; then
                echo "❌ Missing directory: $dir"
                missing=true
            else
                echo "✅ Found directory: $dir"
            fi
        done
        
        # Check key files
        local files=(
            "pom.xml"
            "README.md"
            "Dockerfile"
            "compose.yaml"
        )
        
        for file in "${files[@]}"; do
            if [ ! -f "$project/$file" ]; then
                echo "❌ Missing file: $file"
                missing=true
            else
                echo "✅ Found file: $file"
            fi
        done
        
        if [ "$missing" = true ]; then
            return 1
        fi
    else
        # Basic template checks
        local dirs=(
            "src/main/java"
            "src/main/resources"
            "src/test/java"
        )
        
        local missing=false
        for dir in "${dirs[@]}"; do
            if [ ! -d "$project/$dir" ]; then
                echo "❌ Missing directory: $dir"
                missing=true
            else
                echo "✅ Found directory: $dir"
            fi
        done
        
        # Check key files
        if [ ! -f "$project/pom.xml" ]; then
            echo "❌ Missing file: pom.xml"
            missing=true
        else
            echo "✅ Found file: pom.xml"
        fi
        
        if [ "$missing" = true ]; then
            return 1
        fi
    fi
    
    echo "✅ Project structure verified successfully."
    return 0
}

# Fix project structure
fix_project_structure() {
    local project="$1"
    local template="$2"
    
    echo "Fixing project structure..."
    
    if [ "$template" = "aws-temporal-auth0" ]; then
        # Create all required directories
        mkdir -p "$project/.github/workflows"
        mkdir -p "$project/.mvn"
        mkdir -p "$project/charts/$project/templates"
        mkdir -p "$project/src/main/java/com/example/$project/config"
        mkdir -p "$project/src/main/java/com/example/$project/controller"
        mkdir -p "$project/src/main/java/com/example/$project/domain/dto"
        mkdir -p "$project/src/main/java/com/example/$project/domain/entity"
        mkdir -p "$project/src/main/java/com/example/$project/exception"
        mkdir -p "$project/src/main/java/com/example/$project/middleware"
        mkdir -p "$project/src/main/java/com/example/$project/messaging"
        mkdir -p "$project/src/main/java/com/example/$project/repository"
        mkdir -p "$project/src/main/java/com/example/$project/service"
        mkdir -p "$project/src/main/java/com/example/$project/temporal/activity/impl"
        mkdir -p "$project/src/main/java/com/example/$project/temporal/worker"
        mkdir -p "$project/src/main/java/com/example/$project/temporal/workflow/impl"
        mkdir -p "$project/src/main/java/com/example/$project/util"
        mkdir -p "$project/src/main/resources/db/migration"
        mkdir -p "$project/src/main/swagger"
        mkdir -p "$project/src/test/java/com/example/$project"

        # Create key files if missing
        if [ ! -f "$project/README.md" ]; then
            echo "# $project" > "$project/README.md"
            echo "" >> "$project/README.md"
            echo "A Spring Boot project created with SpringWell CLI." >> "$project/README.md"
            echo "" >> "$project/README.md"
            echo "## Features" >> "$project/README.md"
            echo "" >> "$project/README.md"
            echo "- AWS Integration" >> "$project/README.md"
            echo "- Temporal Workflow Engine" >> "$project/README.md"
            echo "- Auth0 Authentication" >> "$project/README.md"
            echo "- Kubernetes Helm Charts" >> "$project/README.md"
            echo "- Docker Compose for local development" >> "$project/README.md"
        fi
        
        if [ ! -f "$project/Dockerfile" ]; then
            cat > "$project/Dockerfile" << EOF
# Multi-stage build Dockerfile for SpringBoot application
FROM eclipse-temurin:17-jdk as build
WORKDIR /app
COPY . .
RUN ./mvnw clean package -DskipTests

FROM eclipse-temurin:17-jre
WORKDIR /app
COPY --from=build /app/target/*.jar app.jar
ENTRYPOINT ["java", "-jar", "app.jar"]
EOF
        fi
        
        if [ ! -f "$project/compose.yaml" ]; then
            cat > "$project/compose.yaml" << EOF
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SPRING_PROFILES_ACTIVE=dev
    depends_on:
      - postgres
      - temporal

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: ${project}
    ports:
      - "5432:5432"
    volumes:
      - postgres-data:/var/lib/postgresql/data

  temporal:
    image: temporalio/auto-setup:1.22.0
    environment:
      - DB=postgresql
      - DB_PORT=5432
      - POSTGRES_USER=postgres
      - POSTGRES_PWD=postgres
      - POSTGRES_SEEDS=postgres
    ports:
      - "7233:7233"
    depends_on:
      - postgres

volumes:
  postgres-data:
EOF
        fi
        
        # Add placeholder files for key Java classes if they don't exist
        local config_dir="$project/src/main/java/com/example/$project/config"
        if [ ! -f "$config_dir/TemporalConfig.java" ]; then
            mkdir -p "$config_dir"
            cat > "$config_dir/TemporalConfig.java" << EOF
package com.example.$project.config;

import io.temporal.client.WorkflowClient;
import io.temporal.serviceclient.WorkflowServiceStubs;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class TemporalConfig {

    @Value("\${temporal.service.address:localhost:7233}")
    private String temporalServiceAddress;

    @Value("\${temporal.namespace:default}")
    private String namespace;

    @Bean
    public WorkflowServiceStubs workflowServiceStubs() {
        return WorkflowServiceStubs.newServiceStubs(
                WorkflowServiceStubs.newConnectivityOptions()
                        .setTargetEndpoint(temporalServiceAddress)
        );
    }

    @Bean
    public WorkflowClient workflowClient(WorkflowServiceStubs workflowServiceStubs) {
        return WorkflowClient.newInstance(workflowServiceStubs,
                WorkflowClient.newOptions().setNamespace(namespace));
    }
}
EOF
        fi
        
        # Create application.yml if it doesn't exist
        local resources_dir="$project/src/main/resources"
        if [ ! -f "$resources_dir/application.yml" ]; then
            mkdir -p "$resources_dir"
            cat > "$resources_dir/application.yml" << EOF
spring:
  application:
    name: $project
  profiles:
    active: local
  datasource:
    url: jdbc:postgresql://localhost:5432/$project
    username: postgres
    password: postgres
    driver-class-name: org.postgresql.Driver
  jpa:
    hibernate:
      ddl-auto: validate
    properties:
      hibernate:
        dialect: org.hibernate.dialect.PostgreSQLDialect
  flyway:
    locations: classpath:db/migration
    baseline-on-migrate: true

server:
  port: 8080
  servlet:
    context-path: /api

# Temporal configuration
temporal:
  service:
    address: localhost:7233
  namespace: default
  taskqueue:
    default: $project-task-queue

# AWS Configuration
aws:
  region: us-west-2
  
# Auth0 Configuration
auth0:
  audience: https://api.$project.com
  domain: $project.us.auth0.com
  
# Logging configuration
logging:
  level:
    root: INFO
    com.example.$project: DEBUG
    io.temporal: INFO
EOF
        fi
    else
        # Basic template
        mkdir -p "$project/src/main/java/com/example/$project"
        mkdir -p "$project/src/main/resources"
        mkdir -p "$project/src/test/java/com/example/$project"
        
        if [ ! -f "$project/README.md" ]; then
            echo "# $project" > "$project/README.md"
            echo "" >> "$project/README.md"
            echo "A Spring Boot project created with SpringWell CLI." >> "$project/README.md"
        fi
        
        # Create application.properties if it doesn't exist
        local resources_dir="$project/src/main/resources"
        if [ ! -f "$resources_dir/application.properties" ]; then
            mkdir -p "$resources_dir"
            cat > "$resources_dir/application.properties" << EOF
# Application properties
spring.application.name=$project
server.port=8080

# Database configuration
spring.datasource.url=jdbc:h2:mem:testdb
spring.datasource.driverClassName=org.h2.Driver
spring.datasource.username=sa
spring.datasource.password=password
spring.jpa.database-platform=org.hibernate.dialect.H2Dialect
spring.h2.console.enabled=true

# Logging configuration
logging.level.root=INFO
logging.level.com.example.$project=DEBUG
EOF
        fi
    fi
    
    echo "✅ Project structure fixed."
}

# Main execution starts here

# Check if the CLI exists
check_cli
if [ $? -ne 0 ]; then
    exit 1
fi

# Verify template files if enabled
if [ "$VERIFY_TEMPLATES" = true ]; then
    verify_templates "$TEMPLATE"
    if [ $? -ne 0 ]; then
        echo "Template verification failed. Cannot proceed with project creation."
        echo "Please ensure that all template files exist before creating a project."
        exit 1
    fi
fi

# Create the project
echo "Creating new SpringWell project: $PROJECT_NAME"
echo "  Template: $TEMPLATE"
echo "  Database: $DB"
echo ""

if [ "$VERBOSE" = true ]; then
    "$CLI_PATH" new "$PROJECT_NAME" --template "$TEMPLATE" --db "$DB" --verbose
else
    "$CLI_PATH" new "$PROJECT_NAME" --template "$TEMPLATE" --db "$DB"
fi

if [ $? -ne 0 ]; then
    echo "Error: Failed to create project"
    exit 1
fi

echo ""
echo "Project created successfully!"

# Verify the project structure
verify_project_structure "$PROJECT_NAME" "$TEMPLATE"
structure_ok=$?

# Fix structure if needed and requested
if [ $structure_ok -ne 0 ] && [ "$FIX_STRUCTURE" = true ]; then
    echo "Some directories are missing. Fixing project structure..."
    fix_project_structure "$PROJECT_NAME" "$TEMPLATE"
    
    # Verify again after fixing
    verify_project_structure "$PROJECT_NAME" "$TEMPLATE"
    structure_ok=$?
    
    if [ $structure_ok -ne 0 ]; then
        echo "Warning: Project structure is still incomplete after fixing. There might be deeper issues."
    else
        echo "✅ Project structure has been successfully fixed!"
    fi
elif [ $structure_ok -ne 0 ]; then
    echo "Warning: Some project directories are missing."
    echo "You can fix this by running this script with the --fix option."
fi

echo "Navigating to project directory..."

# Navigate to the project directory
cd "$PROJECT_NAME" || { 
    echo "Error: Failed to navigate to project directory"
    exit 1
}

echo "You are now in the $PROJECT_NAME directory"
echo "Ready to start development!"

# Execute a new shell in the project directory
exec $SHELL 
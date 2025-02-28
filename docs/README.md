# SpringWell CLI Documentation

## Table of Contents

1. [Installation](#installation)
2. [Core Commands](#core-commands)
3. [Generation Commands](#generation-commands)
4. [Configuration](#configuration)
5. [Templates](#templates)
6. [Best Practices](#best-practices)

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/springwell/cli.git

# Navigate to the directory
cd cli

# Build the binary
go build -o springwell ./cmd/springwell

# Move the binary to a directory in your PATH
sudo mv springwell /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/springwell/cli/cmd/springwell@latest
```

### Using the Install Script

```bash
# Clone the repository
git clone https://github.com/springwell/cli.git

# Navigate to the directory
cd cli

# Run the install script
./install.sh
```

## Core Commands

### Creating a New Project

```bash
# Create a new project with interactive prompts
springwell new my-service

# Create with specific options
springwell new my-service --package com.company.service --db postgres --auth jwt
```

Options:
- `--package, -p <package>`: Java package name (default: derived from name)
- `--db <database>`: Database type (postgres, mysql, h2) (default: postgres)
- `--auth <type>`: Authentication type (jwt, oauth2, basic) (default: jwt)
- `--cloud <provider>`: Cloud provider integration (aws, azure, gcp) (default: aws)
- `--features <list>`: Comma-separated list of features to include

### Running in Development Mode

```bash
# Run with default settings
springwell dev

# Run with custom port and profile
springwell dev --port 9000 --profile local --debug
```

Options:
- `--port <port>`: Port to run on (default: 8080)
- `--profile <profile>`: Spring profile to use (default: dev)
- `--debug`: Enable debug logging

### Building the Application

```bash
springwell build
```

### Running Tests

```bash
# Run all tests
springwell test

# Run a specific test
springwell test --test UserServiceTest
```

Options:
- `--test <test>`: Specific test to run

### Checking Project Health

```bash
springwell doctor
```

## Generation Commands

### Generating an Entity

```bash
# Generate an entity with interactive field definition
springwell generate entity Product

# Generate with fields specified
springwell generate entity Product --fields "name:String price:Double quantity:Integer description:String:nullable"

# Generate with relationships
springwell generate entity Order --fields "orderDate:Date status:String" --relations "manyToOne:customer:User oneToMany:items:OrderItem"
```

Options:
- `--fields, -f <fields>`: Field definitions (format: "name:type[:modifier]")
- `--relations, -r <relations>`: Relationship definitions (format: "type:field:entity")
- `--table, -t <name>`: Database table name (default: derived from entity name)
- `--audit`: Add auditing fields (created/updated timestamps)
- `--lombok`: Use Lombok annotations
- `--dto`: Generate DTO classes
- `--no-repository`: Skip repository generation
- `--no-service`: Skip service generation
- `--no-controller`: Skip controller generation

### Generating a Controller

```bash
springwell generate controller User
```

### Generating a Service

```bash
springwell generate service User
```

### Generating a Repository

```bash
springwell generate repository User
```

### Generating a DTO

```bash
springwell generate dto User
```

## Configuration

SpringWell can be configured via a `.springwell.yml` file in your project root:

```yaml
project:
  package: com.acme.service
  defaultsDirectory: .springwell/templates

code:
  style:
    indentation: 4
    lineWidth: 120
  lombok: true
  standardizeFields: true
  
templates:
  directory: .springwell/templates
  
aws:
  region: us-east-1
  defaultServices:
    - s3
    - secretsManager
```

## Templates

SpringWell uses templates to generate code. You can customize these templates by creating a `.springwell/templates` directory in your project and copying the default templates there.

The default templates are located in the `pkg/templates` directory of the SpringWell CLI repository.

## Best Practices

1. **Consistent Naming**: Use consistent naming conventions for your entities, services, and controllers.
2. **Field Definitions**: When defining fields, use the format `name:type[:modifier]` where `modifier` can be `nullable`.
3. **Relationship Definitions**: When defining relationships, use the format `type:field:entity` where `type` can be `oneToOne`, `oneToMany`, `manyToOne`, or `manyToMany`.
4. **Custom Templates**: Create custom templates to match your project's coding style and standards.
5. **Project Structure**: Follow the standard Spring Boot project structure for better maintainability. 
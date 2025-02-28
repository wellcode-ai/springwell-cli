# SpringWell CLI

SpringWell is a developer-friendly CLI tool designed to standardize and accelerate development with the Spring Boot Enterprise Starter Template. It helps you generate consistent code components, manage your project, and enforce best practices across your team.

# Enterprise-Ready Spring Boot Template

## Accelerate Your AWS Cloud Journey

This production-grade Spring Boot template empowers your team to build scalable, resilient microservices in record time. Engineered with enterprise standards in mind, it seamlessly integrates the most powerful tools in modern cloud architecture.

### 🚀 Key Features

- **Battle-Tested Architecture**: Built on Spring Boot 3.x with Java 17, following industry best practices for enterprise applications
- **Temporal Workflow Engine**: Implement resilient, distributed workflows that survive infrastructure failures
- **Auth0 Integration**: Enterprise-grade security with minimal configuration
- **AWS-Optimized**: First-class integration with AWS services (S3, SQS, Secrets Manager, RDS)
- **Datadog Observability**: Comprehensive monitoring, tracing, and analytics out of the box
- **Distributed Caching**: Redis/ElastiCache integration for high-performance data access
- **Infrastructure as Code**: Complete CloudFormation templates for reliable deployments
- **CI/CD Ready**: GitHub Actions workflows for automated testing and deployment

### 💼 Business Benefits

- **Reduce Time-to-Market**: Start building features on day one instead of spending weeks on infrastructure
- **Lower Operational Costs**: Pre-configured observability and best practices reduce production incidents
- **Decrease Onboarding Time**: Standardized project structure makes it easy for new team members to contribute
- **Enterprise Compliance**: Security-first approach with audit trails and proper authentication
- **Future-Proof Investment**: Regular updates to keep dependencies secure and current

### 🛠️ Technical Excellence

This template represents hundreds of hours of engineering expertise, distilled into a clean, modular foundation for your next project. Each component has been carefully selected and integrated to create a cohesive development experience.

Start building what matters today, and leave the infrastructure complexities to us.

---

*"This template saved our team at least 6 weeks of setup and configuration time. We were able to focus on business logic from day one."* - Engineering Director at Enterprise Customer

## Installation

### From Source
```bash
# Clone the repository
git clone https://github.com/springwell/cli.git
cd cli

# Build the CLI
go build -o springwell cmd/springwell/main.go

# Make it executable
chmod +x springwell

# Optional: Move to a directory in your PATH
sudo mv springwell /usr/local/bin/
```

### Using Go Install
```bash
go install github.com/springwell/cli/cmd/springwell@latest
```

### Using Install Script
```bash
curl -sSL https://get.springwell.dev | bash
```

## Usage

SpringWell CLI has two modes of operation: command-line mode and interactive mode.

### Command-Line Mode

```bash
# Create a new Spring Boot project
springwell new my-project

# Create a project with a specific template
springwell new my-project --template aws-temporal-auth0

# Run application in development mode
springwell dev

# Run with custom port and profile
springwell dev --port 8081 --profile local

# Build the application
springwell build

# Run tests
springwell test

# Run a specific test
springwell test --test UserServiceTest

# Check project health
springwell doctor

# Generate components
springwell generate entity User name:String email:String:nullable

# Generate a controller
springwell generate controller User

# Generate a service
springwell generate service User

# Get help with any command
springwell help generate
```

### Interactive Mode

SpringWell CLI offers an interactive mode that guides you through various operations with a menu-based interface. This is especially useful for new users or when exploring features.

```bash
# Start the CLI in interactive mode
springwell interactive
# or use the shorthand
springwell i
```

In interactive mode, you can:

1. Generate a brand new project with guided prompts
2. Generate various components (entities, controllers, services, etc.)
3. Run your application in development mode
4. Build your project
5. Run tests
6. Check project health

The interactive mode will walk you through each operation step by step, prompting for required information.

### Common Options

All commands support the following options:

- `--quiet` or `-q`: Suppress non-error output
- `--json`: Output in JSON format (for scripting)
- `--help`: Show help information
- `--version`: Show version information

## Examples

### Creating a New Project with AWS, Temporal, and Auth0 Integration

```bash
springwell new my-aws-app --template aws-temporal-auth0 --db postgres
```

### Generating Complete API Endpoints

```bash
# Generate an entity
springwell generate entity Product name:String price:Double description:String:nullable

# Generate a repository
springwell generate repository Product

# Generate a service
springwell generate service Product

# Generate a controller
springwell generate controller Product
```

### Generating Temporal Workflows

```bash
# Generate a Temporal workflow + activity
springwell generate workflow OrderProcessing
```

## Configuration

SpringWell CLI uses a `.springwell.yml` file in your project root to store configuration. This file is automatically created when you create a new project.

Example configuration:

```yaml
project:
  name: my-project
  package: com.example.myproject
  template: aws-temporal-auth0
  database: postgres
```

## Templates

SpringWell CLI comes with pre-defined templates:

- `basic`: Standard Spring Boot project
- `aws-temporal-auth0`: AWS-optimized Spring Boot with Temporal workflows and Auth0 authentication

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
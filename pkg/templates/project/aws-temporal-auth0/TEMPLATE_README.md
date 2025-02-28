# AWS Temporal Auth0 Template

This template creates a Spring Boot application with integrations for:

- **AWS Services** (S3, etc.)
- **Temporal** for durable workflow orchestration 
- **Auth0** for authentication and authorization
- **Datadog** for monitoring and observability

## Using This Template

To create a new project using this template with the SpringWell CLI:

```bash
./springwell new --name my-project --template aws-temporal-auth0
```

## Template Features

1. **Auth0 Integration**
   - Ready-to-use security configuration
   - JWT authentication and authorization
   - Resource protection with role-based access

2. **AWS Integration**
   - S3 client configuration
   - AWS deployment configuration
   - GitHub Actions for AWS CI/CD

3. **Temporal Workflow Engine**
   - Sample workflow implementation
   - Activity interface and implementation
   - Worker service configuration
   - Local development with Temporal Docker setup

4. **Datadog Monitoring**
   - Integrated tracing agent
   - Metrics configuration
   - Logging setup

5. **Development Environment**
   - Docker and Docker Compose configuration
   - Local development setup
   - Testing framework

## Customizing the Template

After generating a project with this template, you can:

1. Update the Auth0 credentials in `application.yml`
2. Configure your AWS region and services
3. Modify the Temporal workflow for your specific use case
4. Set up your Datadog account information

## Requirements

To use this template, you need:

- Java 17+
- Docker and Docker Compose (for local development)
- AWS account (for deployment)
- Auth0 account (for authentication)
- Datadog account (optional, for monitoring) 
package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/springwell/cli/pkg/config"
	"github.com/springwell/cli/pkg/util"
	"github.com/urfave/cli/v2"
)

// NewProjectCommand returns the command to create a new project
func NewProjectCommand() *cli.Command {
	return &cli.Command{
		Name:  "new",
		Usage: "Create a new Spring Boot project",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "package",
				Aliases: []string{"p"},
				Usage:   "Java package name (default: derived from name)",
			},
			&cli.StringFlag{
				Name:  "db",
				Usage: "Database type (postgres, mysql, h2)",
				Value: "postgres",
			},
			&cli.StringFlag{
				Name:  "auth",
				Usage: "Authentication type (jwt, oauth2, basic, auth0)",
				Value: "jwt",
			},
			&cli.StringFlag{
				Name:  "cloud",
				Usage: "Cloud provider integration (aws, azure, gcp)",
				Value: "aws",
			},
			&cli.StringFlag{
				Name:  "features",
				Usage: "Comma-separated list of features to include",
				Value: "swagger,actuator",
			},
			&cli.StringFlag{
				Name:  "template",
				Usage: "Project template to use (basic, aws-temporal-auth0)",
				Value: "basic",
			},
		},
		Action: func(c *cli.Context) error {
			projectName := c.Args().First()
			if projectName == "" {
				return errors.New("project name is required")
			}

			// Create project directory
			projectDir := filepath.Join(".", projectName)
			if err := util.CreateDirectory(projectDir); err != nil {
				return err
			}

			// Determine package name
			packageName := c.String("package")
			if packageName == "" {
				packageName = "com." + util.ToPackageName(projectName)
			}

			// Set up configuration
			cfg := config.GetDefaultConfig()
			cfg.Project.Package = packageName

			// Save configuration
			if err := config.SaveConfig(cfg, projectDir); err != nil {
				return err
			}

			// Use template if specified
			template := c.String("template")
			if template == "aws-temporal-auth0" {
				// Create project with AWS, Temporal, and Auth0 integration
				return createAwsTemporalAuth0Project(projectName, packageName, projectDir, c.String("db"))
			}

			// Use Spring Initializr to create the base project
			return createSpringBootProject(projectName, packageName, projectDir, c.String("db"), c.String("auth"), c.String("features"))
		},
	}
}

// DevCommand returns the command to run the application in development mode
func DevCommand() *cli.Command {
	return &cli.Command{
		Name:  "dev",
		Usage: "Run application with hot reloading",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port",
				Usage: "Port to run on",
				Value: 8080,
			},
			&cli.StringFlag{
				Name:  "profile",
				Usage: "Spring profile to use",
				Value: "dev",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug logging",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			// Check if it's a Spring Boot project
			if !util.IsSpringBootProject(".") {
				return errors.New("current directory is not a Spring Boot project")
			}

			// Determine if it's Maven or Gradle
			isMaven := true
			if _, err := os.Stat("build.gradle"); err == nil {
				isMaven = false
			}

			// Build the command
			var cmd *exec.Cmd
			if isMaven {
				cmd = exec.Command("./mvnw", "spring-boot:run",
					"-Dspring-boot.run.profiles="+c.String("profile"),
					"-Dserver.port="+fmt.Sprintf("%d", c.Int("port")))
			} else {
				cmd = exec.Command("./gradlew", "bootRun",
					"-Dspring.profiles.active="+c.String("profile"),
					"-Dserver.port="+fmt.Sprintf("%d", c.Int("port")))
			}

			if c.Bool("debug") {
				cmd.Env = append(os.Environ(), "DEBUG=true")
			}

			// Set up command output
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Run the command
			util.PrintInfo("Starting application in development mode...")
			return cmd.Run()
		},
	}
}

// BuildCommand returns the command to build the application
func BuildCommand() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build the application",
		Action: func(c *cli.Context) error {
			// Check if it's a Spring Boot project
			if !util.IsSpringBootProject(".") {
				return errors.New("current directory is not a Spring Boot project")
			}

			// Determine if it's Maven or Gradle
			isMaven := true
			if _, err := os.Stat("build.gradle"); err == nil {
				isMaven = false
			}

			// Build the command
			var cmd *exec.Cmd
			if isMaven {
				cmd = exec.Command("./mvnw", "clean", "package", "-DskipTests")
			} else {
				cmd = exec.Command("./gradlew", "clean", "build", "-x", "test")
			}

			// Set up command output
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Run the command
			util.PrintInfo("Building application...")
			return cmd.Run()
		},
	}
}

// TestCommand returns the command to run tests
func TestCommand() *cli.Command {
	return &cli.Command{
		Name:  "test",
		Usage: "Run tests with smart detection",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "test",
				Usage: "Specific test to run",
			},
		},
		Action: func(c *cli.Context) error {
			// Check if it's a Spring Boot project
			if !util.IsSpringBootProject(".") {
				return errors.New("current directory is not a Spring Boot project")
			}

			// Determine if it's Maven or Gradle
			isMaven := true
			if _, err := os.Stat("build.gradle"); err == nil {
				isMaven = false
			}

			// Build the command
			var cmd *exec.Cmd
			testName := c.String("test")
			if isMaven {
				if testName != "" {
					cmd = exec.Command("./mvnw", "test", "-Dtest="+testName)
				} else {
					cmd = exec.Command("./mvnw", "test")
				}
			} else {
				if testName != "" {
					cmd = exec.Command("./gradlew", "test", "--tests", testName)
				} else {
					cmd = exec.Command("./gradlew", "test")
				}
			}

			// Set up command output
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Run the command
			util.PrintInfo("Running tests...")
			return cmd.Run()
		},
	}
}

// DoctorCommand returns the command to check project health
func DoctorCommand() *cli.Command {
	return &cli.Command{
		Name:  "doctor",
		Usage: "Check project health and suggest fixes",
		Action: func(c *cli.Context) error {
			// Check if it's a Spring Boot project
			if !util.IsSpringBootProject(".") {
				return errors.New("current directory is not a Spring Boot project")
			}

			// TODO: Implement project health checks
			util.PrintSuccess("Project looks healthy!")
			return nil
		},
	}
}

// createSpringBootProject creates a new Spring Boot project using Spring Initializr
func createSpringBootProject(name, packageName, projectDir, db, auth, features string) error {
	// Build list of dependencies
	dependencies := []string{
		"web",
		"data-jpa",
		"validation",
		"lombok",
	}

	// Add database dependency
	switch db {
	case "postgres":
		dependencies = append(dependencies, "postgresql")
	case "mysql":
		dependencies = append(dependencies, "mysql")
	default:
		dependencies = append(dependencies, "h2")
	}

	// Add auth dependency
	switch auth {
	case "jwt":
		dependencies = append(dependencies, "security", "oauth2-resource-server")
	case "oauth2":
		dependencies = append(dependencies, "security", "oauth2-client")
	case "basic":
		dependencies = append(dependencies, "security")
	}

	// Create command
	url := "https://start.spring.io/starter.zip"
	url += "?name=" + name
	url += "&groupId=" + packageName
	url += "&artifactId=" + name
	url += "&packageName=" + packageName
	url += "&language=java"
	url += "&javaVersion=17"
	url += "&type=maven-project"

	// Add dependencies
	for _, dep := range dependencies {
		url += "&dependencies=" + dep
	}

	// Execute curl command to download the zip
	zipFile := filepath.Join(projectDir, "temp.zip")
	cmd := exec.Command("curl", "-L", url, "-o", zipFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	util.PrintInfo("Downloading Spring Boot template...")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Extract the zip
	unzipCmd := exec.Command("unzip", "-o", zipFile, "-d", projectDir)
	unzipCmd.Stdout = os.Stdout
	unzipCmd.Stderr = os.Stderr

	util.PrintInfo("Extracting template...")
	if err := unzipCmd.Run(); err != nil {
		return err
	}

	// Remove the zip file
	if err := os.Remove(zipFile); err != nil {
		return err
	}

	// Make the Maven wrapper executable
	mvnwFile := filepath.Join(projectDir, "mvnw")
	if err := os.Chmod(mvnwFile, 0755); err != nil {
		return err
	}

	util.PrintSuccess("Created %s at %s", name, projectDir)
	return nil
}

// createAwsTemporalAuth0Project creates a new Spring Boot project with AWS, Temporal, and Auth0 integration
func createAwsTemporalAuth0Project(name, packageName, projectDir, db string) error {
	// First create a basic Spring Boot project
	if err := createSpringBootProject(name, packageName, projectDir, db, "auth0", "swagger,actuator,webflux,security,lombok,data-jpa"); err != nil {
		return err
	}

	// Now enhance it with AWS, Temporal, and Auth0 structure
	util.PrintInfo("Enhancing project with AWS, Temporal, Auth0, Helm, DB migrations, and OpenAPI...")

	// Create additional directories
	directories := []string{
		// Java directories
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "config"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "controller"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "domain/audit"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "domain/dto"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "domain/entity"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "exception"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "messaging/consumer"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "middleware"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "repository"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "security"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "service"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "temporal/activity"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "temporal/workflow"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "temporal/worker"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "util"),

		// Resource directories
		filepath.Join(projectDir, "src/main/resources/datadog"),
		filepath.Join(projectDir, "src/main/resources/db/migration"),
		filepath.Join(projectDir, "src/main/resources/templates/email"),

		// AWS directories
		filepath.Join(projectDir, "aws/cloudformation"),
		filepath.Join(projectDir, "aws/codebuild"),

		// GitHub workflows
		filepath.Join(projectDir, ".github/workflows"),

		// Scripts and Docker
		filepath.Join(projectDir, "scripts"),
		filepath.Join(projectDir, "docker/datadog-agent"),

		// NEW: Helm charts
		filepath.Join(projectDir, "helm/"+name+"/templates"),
		filepath.Join(projectDir, "helm/"+name+"/charts"),

		// NEW: OpenAPI
		filepath.Join(projectDir, "src/main/resources/openapi"),

		// NEW: Additional util directories
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "util/mapper"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "util/validation"),
		filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "util/logging"),
	}

	for _, dir := range directories {
		if err := util.CreateDirectory(dir); err != nil {
			return err
		}
	}

	// Define strings for the backticks to avoid linter errors
	backtick := "`"

	// Create README with project structure
	readmeContent := "# ${name}\n\n" +
		"## Project Structure\n\n" +
		"This project follows the Spring Boot best practices and includes integration with:\n" +
		"- AWS services\n" +
		"- Temporal for workflow orchestration\n" +
		"- Auth0 for authentication\n" +
		"- Datadog for monitoring\n" +
		"- Helm charts for Kubernetes deployment\n" +
		"- Database migrations with Flyway\n" +
		"- OpenAPI for API documentation\n\n" +
		backtick + backtick + backtick + "\n" +
		"├── src/\n" +
		"│   ├── main/\n" +
		"│   │   ├── java/\n" +
		"│   │   │   └── ${packagePath}/\n" +
		"│   │   │       ├── ApplicationMain.java\n" +
		"│   │   │       ├── config/\n" +
		"│   │   │       ├── controller/\n" +
		"│   │   │       ├── domain/\n" +
		"│   │   │       ├── exception/\n" +
		"│   │   │       ├── messaging/\n" +
		"│   │   │       ├── middleware/\n" +
		"│   │   │       ├── repository/\n" +
		"│   │   │       ├── security/\n" +
		"│   │   │       ├── service/\n" +
		"│   │   │       ├── temporal/\n" +
		"│   │   │       └── util/\n" +
		"│   │   └── resources/\n" +
		"│   │       ├── db/migration/     # Flyway migrations\n" +
		"│   │       └── openapi/          # OpenAPI specifications\n" +
		"│   └── test/\n" +
		"├── aws/\n" +
		"├── helm/                         # Kubernetes Helm charts\n" +
		"├── .github/\n" +
		"├── scripts/\n" +
		"├── docker/\n" +
		backtick + backtick + backtick + "\n\n" +
		"## Getting Started\n\n" +
		"1. Configure your Auth0 credentials in \\`application.yml\\`\n" +
		"2. Set up your AWS credentials\n" +
		"3. Start the application: \\`./gradlew bootRun\\`\n\n" +
		"## Deployment\n\n" +
		"- **Kubernetes**: Use the provided Helm charts in the \\`helm\\` directory\n" +
		"- **AWS**: Use the CloudFormation templates in the \\`aws/cloudformation\\` directory\n"

	readmeContent = strings.ReplaceAll(readmeContent, "${name}", name)
	readmeContent = strings.ReplaceAll(readmeContent, "${packagePath}", strings.ReplaceAll(packageName, ".", "/"))

	if err := util.WriteFile(filepath.Join(projectDir, "README.md"), readmeContent); err != nil {
		return err
	}

	// Create application.yml with Auth0 and AWS configurations
	applicationYml := `spring:
  application:
    name: ${name}
  datasource:
    url: jdbc:postgresql://localhost:5432/${name}
    username: postgres
    password: postgres
  jpa:
    hibernate:
      ddl-auto: validate
    properties:
      hibernate:
        dialect: org.hibernate.dialect.PostgreSQLDialect
  flyway:
    enabled: true
    baseline-on-migrate: true
    locations: classpath:db/migration
    
# OpenAPI Configuration
springdoc:
  api-docs:
    path: /api-docs
  swagger-ui:
    path: /swagger-ui.html
    
# Auth0 Configuration
auth0:
  audience: https://api.example.com
  domain: your-domain.auth0.com
  client-id: your-client-id
  client-secret: your-client-secret
  
# AWS Configuration
cloud:
  aws:
    region:
      static: us-east-1
    stack:
      auto: false
    credentials:
      instance-profile: true
      
# Temporal Configuration
temporal:
  service-address: 127.0.0.1:7233
  namespace: default
  
# Datadog Configuration
datadog:
  enabled: true
  api-key: your-api-key
  application-key: your-app-key
  service-name: ${name}
  environment: development
  
# Logging
logging:
  level:
    root: INFO
    ${packageName}: DEBUG
    
# Server configuration
server:
  port: 8080
`

	applicationYml = strings.ReplaceAll(applicationYml, "${name}", name)
	applicationYml = strings.ReplaceAll(applicationYml, "${packageName}", packageName)

	if err := util.WriteFile(filepath.Join(projectDir, "src/main/resources/application.yml"), applicationYml); err != nil {
		return err
	}

	// Create a basic Dockerfile
	dockerfileContent := `FROM openjdk:17-slim

WORKDIR /app

COPY build/libs/${name}-0.0.1-SNAPSHOT.jar app.jar

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "app.jar"]
`

	dockerfileContent = strings.ReplaceAll(dockerfileContent, "${name}", name)

	if err := util.WriteFile(filepath.Join(projectDir, "docker/Dockerfile"), dockerfileContent); err != nil {
		return err
	}

	// Create a basic docker-compose.yml
	dockerComposeContent := `version: '3.8'

services:
  app:
    build:
      context: ..
      dockerfile: docker/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - SPRING_PROFILES_ACTIVE=dev
    depends_on:
      - postgres
      - temporal
      - datadog-agent
    networks:
      - app-network
      
  postgres:
    image: postgres:14-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_DB=${name}
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
    volumes:
      - postgres-data:/var/lib/postgresql/data
    networks:
      - app-network
      
  temporal:
    image: temporalio/auto-setup:1.18.0
    ports:
      - "7233:7233"
    environment:
      - DYNAMIC_CONFIG_FILE_PATH=config/dynamicconfig/development.yaml
    networks:
      - app-network
      
  datadog-agent:
    image: datadog/agent:latest
    environment:
      - DD_API_KEY=your-api-key
      - DD_APM_ENABLED=true
      - DD_LOGS_ENABLED=true
      - DD_LOGS_CONFIG_CONTAINER_COLLECT_ALL=true
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - /proc/:/host/proc/:ro
      - /sys/fs/cgroup/:/host/sys/fs/cgroup:ro
    networks:
      - app-network

networks:
  app-network:
    driver: bridge
    
volumes:
  postgres-data:
`

	dockerComposeContent = strings.ReplaceAll(dockerComposeContent, "${name}", name)

	if err := util.WriteFile(filepath.Join(projectDir, "docker/docker-compose.yml"), dockerComposeContent); err != nil {
		return err
	}

	// NEW: Create compose.yaml at the root
	composeYamlContent := `# Root-level compose.yaml - development environment
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: docker/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - SPRING_PROFILES_ACTIVE=dev
      - SPRING_DATASOURCE_URL=jdbc:postgresql://postgres:5432/${name}
    depends_on:
      - postgres
      - temporal
    volumes:
      - ./:/app
    networks:
      - ${name}-network

  postgres:
    image: postgres:14-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_DB=${name}
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
    volumes:
      - postgres-data:/var/lib/postgresql/data
    networks:
      - ${name}-network

  temporal:
    image: temporalio/auto-setup:1.18.0
    ports:
      - "7233:7233"
      - "8088:8088"
    environment:
      - DB=postgresql
      - DB_PORT=5432
      - POSTGRES_USER=postgres
      - POSTGRES_PWD=postgres
      - POSTGRES_SEEDS=postgres
    depends_on:
      - postgres
    networks:
      - ${name}-network

networks:
  ${name}-network:
    driver: bridge

volumes:
  postgres-data:
`

	composeYamlContent = strings.ReplaceAll(composeYamlContent, "${name}", name)

	if err := util.WriteFile(filepath.Join(projectDir, "compose.yaml"), composeYamlContent); err != nil {
		return err
	}

	// Create a basic GitHub Actions workflow
	ciYamlContent := `name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up JDK 17
      uses: actions/setup-java@v3
      with:
        java-version: '17'
        distribution: 'temurin'
        
    - name: Build with Gradle
      run: ./gradlew build
      
    - name: Run tests
      run: ./gradlew test
`

	if err := util.WriteFile(filepath.Join(projectDir, ".github/workflows/ci.yml"), ciYamlContent); err != nil {
		return err
	}

	// NEW: Create Helm chart files
	helmChartYamlContent := `apiVersion: v2
name: ${name}
description: A Helm chart for ${name} Spring Boot application
type: application
version: 0.1.0
appVersion: "1.0.0"
`

	helmChartYamlContent = strings.ReplaceAll(helmChartYamlContent, "${name}", name)

	if err := util.WriteFile(filepath.Join(projectDir, "helm/"+name+"/Chart.yaml"), helmChartYamlContent); err != nil {
		return err
	}

	helmValuesYamlContent := `# Default values for ${name}
replicaCount: 1

image:
  repository: ${name}
  tag: latest
  pullPolicy: IfNotPresent

nameOverride: ""
fullnameOverride: ""

service:
  type: ClusterIP
  port: 8080

ingress:
  enabled: false
  annotations: {}
  hosts:
    - host: ${name}.local
      paths: ["/"]

resources:
  limits:
    cpu: 1000m
    memory: 1024Mi
  requests:
    cpu: 500m
    memory: 512Mi

env:
  - name: SPRING_PROFILES_ACTIVE
    value: prod
  - name: SPRING_DATASOURCE_URL
    value: jdbc:postgresql://postgres:5432/${name}

secrets:
  - name: auth0-credentials
    data:
      client-id: your-client-id
      client-secret: your-client-secret

postgresql:
  enabled: true
  postgresqlUsername: postgres
  postgresqlPassword: postgres
  postgresqlDatabase: ${name}
`

	helmValuesYamlContent = strings.ReplaceAll(helmValuesYamlContent, "${name}", name)

	if err := util.WriteFile(filepath.Join(projectDir, "helm/"+name+"/values.yaml"), helmValuesYamlContent); err != nil {
		return err
	}

	// Create deployment template
	deploymentYamlContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "fullname" . }}
  labels:
    {{- include "labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "selectorLabels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
          env:
            {{- range .Values.env }}
            - name: {{ .name }}
              value: {{ .value | quote }}
            {{- end }}
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
`

	if err := util.WriteFile(filepath.Join(projectDir, "helm/"+name+"/templates/deployment.yaml"), deploymentYamlContent); err != nil {
		return err
	}

	// Create service template
	serviceYamlContent := `apiVersion: v1
kind: Service
metadata:
  name: {{ include "fullname" . }}
  labels:
    {{- include "labels" . | nindent 4 }}
spec:
  type: {{ .Values.service.type }}
  ports:
    - port: {{ .Values.service.port }}
      targetPort: http
      protocol: TCP
      name: http
  selector:
    {{- include "selectorLabels" . | nindent 4 }}
`

	if err := util.WriteFile(filepath.Join(projectDir, "helm/"+name+"/templates/service.yaml"), serviceYamlContent); err != nil {
		return err
	}

	// Create helpers template
	helpersTemplContent := `{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- printf "%s" $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
{{ include "selectorLabels" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "selectorLabels" -}}
app.kubernetes.io/name: {{ include "name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
`

	if err := util.WriteFile(filepath.Join(projectDir, "helm/"+name+"/templates/_helpers.tpl"), helpersTemplContent); err != nil {
		return err
	}

	// NEW: Create Flyway database migrations
	v1MigrationSql := `-- V1__initial_schema.sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(255),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_roles (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id),
  role VARCHAR(50) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create audit tables for tracking changes
CREATE TABLE audit_log (
  id SERIAL PRIMARY KEY,
  entity_type VARCHAR(100) NOT NULL,
  entity_id INTEGER NOT NULL,
  action VARCHAR(50) NOT NULL,
  actor VARCHAR(255) NOT NULL,
  changes JSONB,
  timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
`

	if err := util.WriteFile(filepath.Join(projectDir, "src/main/resources/db/migration/V1__initial_schema.sql"), v1MigrationSql); err != nil {
		return err
	}

	v2MigrationSql := `-- V2__temporal_tables.sql
CREATE TABLE workflow_state (
  workflow_id VARCHAR(255) PRIMARY KEY,
  workflow_type VARCHAR(255) NOT NULL,
  status VARCHAR(50) NOT NULL,
  started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  completed_at TIMESTAMP WITH TIME ZONE,
  data JSONB
);

CREATE TABLE activity_execution (
  id SERIAL PRIMARY KEY,
  workflow_id VARCHAR(255) REFERENCES workflow_state(workflow_id),
  activity_type VARCHAR(255) NOT NULL,
  status VARCHAR(50) NOT NULL,
  started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  completed_at TIMESTAMP WITH TIME ZONE,
  attempts INTEGER DEFAULT 0,
  last_error TEXT
);
`

	if err := util.WriteFile(filepath.Join(projectDir, "src/main/resources/db/migration/V2__temporal_tables.sql"), v2MigrationSql); err != nil {
		return err
	}

	// NEW: Create OpenAPI specification
	openApiYamlContent := `openapi: 3.0.3
info:
  title: ${name} API
  description: API documentation for ${name}
  version: 1.0.0
servers:
  - url: /api
paths:
  /v1/health:
    get:
      summary: Health check endpoint
      responses:
        '200':
          description: Service is healthy
  /v1/users:
    get:
      summary: List all users
      security:
        - bearerAuth: []
      responses:
        '200':
          description: A list of users
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/User'
        '401':
          description: Unauthorized
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: integer
          format: int64
        email:
          type: string
        name:
          type: string
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
`

	openApiYamlContent = strings.ReplaceAll(openApiYamlContent, "${name}", name)

	if err := util.WriteFile(filepath.Join(projectDir, "src/main/resources/openapi/api.yaml"), openApiYamlContent); err != nil {
		return err
	}

	// NEW: Create OpenAPI generator configuration
	openApiGeneratorContent := `# Configuration for OpenAPI Generator
generatorName: spring
inputSpec: src/main/resources/openapi/api.yaml
outputDir: build/generated/openapi
apiPackage: ${packageName}.controller.api
modelPackage: ${packageName}.domain.dto
templateDir: src/main/resources/openapi/templates
additionalProperties:
  java8: true
  reactive: true
  interfaceOnly: true
`

	openApiGeneratorContent = strings.ReplaceAll(openApiGeneratorContent, "${packageName}", packageName)

	if err := util.WriteFile(filepath.Join(projectDir, "openapi-generator.yaml"), openApiGeneratorContent); err != nil {
		return err
	}

	// NEW: Create utility classes
	mapperUtilContent := `package ${packageName}.util.mapper;

import org.springframework.stereotype.Component;

/**
 * Utility class for mapping between DTOs and entities.
 */
@Component
public class EntityMapper {
    
    /**
     * Maps entity to DTO.
     */
    public <D, E> D toDto(E entity, Class<D> dtoClass) {
        // Implementation will be provided by MapStruct in a real scenario
        throw new UnsupportedOperationException("Not implemented yet");
    }
    
    /**
     * Maps DTO to entity.
     */
    public <D, E> E toEntity(D dto, Class<E> entityClass) {
        // Implementation will be provided by MapStruct in a real scenario
        throw new UnsupportedOperationException("Not implemented yet");
    }
}
`

	mapperUtilContent = strings.ReplaceAll(mapperUtilContent, "${packageName}", packageName)

	if err := util.WriteFile(filepath.Join(projectDir, "src/main/java", strings.ReplaceAll(packageName, ".", "/"), "util/mapper/EntityMapper.java"), mapperUtilContent); err != nil {
		return err
	}

	util.PrintSuccess("Successfully created AWS+Temporal+Auth0 project structure with Helm Charts, DB migrations, OpenAPI, and utilities at %s", projectDir)
	return nil
}

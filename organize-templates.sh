#!/bin/bash

# organize-templates.sh - Script to organize template files into the proper directory structure

echo "Organizing template files into the proper directory structure..."

# Base template directory
TEMPLATE_DIR="pkg/templates/project/aws-temporal-auth0"

# Create necessary directories
mkdir -p "$TEMPLATE_DIR/config"
mkdir -p "$TEMPLATE_DIR/temporal/worker"
mkdir -p "$TEMPLATE_DIR/temporal/workflow"
mkdir -p "$TEMPLATE_DIR/temporal/activity"
mkdir -p "$TEMPLATE_DIR/controller"
mkdir -p "$TEMPLATE_DIR/middleware"
mkdir -p "$TEMPLATE_DIR/messaging"
mkdir -p "$TEMPLATE_DIR/exception"
mkdir -p "$TEMPLATE_DIR/domain/dto"
mkdir -p "$TEMPLATE_DIR/domain/entity"
mkdir -p "$TEMPLATE_DIR/repository"
mkdir -p "$TEMPLATE_DIR/service"
mkdir -p "$TEMPLATE_DIR/util"

# Move config files to the config directory if they're in the root
for config_file in Auth0Config.java.tmpl AwsS3Config.java.tmpl DatadogConfig.java.tmpl SecurityConfig.java.tmpl TemporalConfig.java.tmpl; do
    if [ -f "$TEMPLATE_DIR/$config_file" ]; then
        echo "Moving $config_file to config directory"
        mv "$TEMPLATE_DIR/$config_file" "$TEMPLATE_DIR/config/"
    fi
done

# Move temporal files to their respective directories
if [ -f "$TEMPLATE_DIR/SampleWorkflow.java.tmpl" ]; then
    echo "Moving SampleWorkflow.java.tmpl to temporal/workflow directory"
    mv "$TEMPLATE_DIR/SampleWorkflow.java.tmpl" "$TEMPLATE_DIR/temporal/workflow/OrderProcessingWorkflow.java.tmpl"
fi

if [ -f "$TEMPLATE_DIR/SampleWorkflowImpl.java.tmpl" ]; then
    echo "Moving SampleWorkflowImpl.java.tmpl to temporal/workflow directory"
    mv "$TEMPLATE_DIR/SampleWorkflowImpl.java.tmpl" "$TEMPLATE_DIR/temporal/workflow/OrderProcessingWorkflowImpl.java.tmpl"
fi

if [ -f "$TEMPLATE_DIR/SampleActivity.java.tmpl" ]; then
    echo "Moving SampleActivity.java.tmpl to temporal/activity directory"
    mv "$TEMPLATE_DIR/SampleActivity.java.tmpl" "$TEMPLATE_DIR/temporal/activity/OrderProcessingActivity.java.tmpl"
fi

if [ -f "$TEMPLATE_DIR/SampleActivityImpl.java.tmpl" ]; then
    echo "Moving SampleActivityImpl.java.tmpl to temporal/activity directory"
    mv "$TEMPLATE_DIR/SampleActivityImpl.java.tmpl" "$TEMPLATE_DIR/temporal/activity/OrderProcessingActivityImpl.java.tmpl"
fi

if [ -f "$TEMPLATE_DIR/WorkerService.java.tmpl" ]; then
    echo "Moving WorkerService.java.tmpl to temporal/worker directory"
    mv "$TEMPLATE_DIR/WorkerService.java.tmpl" "$TEMPLATE_DIR/temporal/worker/TemporalWorkerService.java.tmpl"
    
    # Create TemporalWorkerRegistrar.java.tmpl if it doesn't exist
    if [ ! -f "$TEMPLATE_DIR/temporal/worker/TemporalWorkerRegistrar.java.tmpl" ]; then
        echo "Creating TemporalWorkerRegistrar.java.tmpl"
        cat > "$TEMPLATE_DIR/temporal/worker/TemporalWorkerRegistrar.java.tmpl" << EOF
package {{package}}.temporal.worker;

import {{package}}.temporal.activity.OrderProcessingActivityImpl;
import {{package}}.temporal.workflow.OrderProcessingWorkflowImpl;
import io.temporal.client.WorkflowClient;
import io.temporal.worker.Worker;
import io.temporal.worker.WorkerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

/**
 * Component responsible for registering Temporal workers, workflows, and activities.
 */
@Component
public class TemporalWorkerRegistrar {

    private final WorkflowClient workflowClient;
    
    @Value("\${temporal.taskqueue.default:{{artifactId}}-task-queue}")
    private String defaultTaskQueue;
    
    @Value("\${temporal.taskqueue.orders:order-processing-queue}")
    private String orderProcessingQueue;
    
    @Value("\${temporal.taskqueue.payments:payment-processing-queue}")
    private String paymentProcessingQueue;
    
    @Value("\${temporal.taskqueue.notifications:notification-queue}")
    private String notificationQueue;

    @Autowired
    public TemporalWorkerRegistrar(WorkflowClient workflowClient) {
        this.workflowClient = workflowClient;
    }

    /**
     * Registers all workers, workflows, and activities.
     * 
     * @param workerFactory The worker factory to create workers
     */
    public void registerWorkers(WorkerFactory workerFactory) {
        // Register order processing workflows and activities
        Worker orderProcessingWorker = workerFactory.newWorker(orderProcessingQueue);
        orderProcessingWorker.registerWorkflowImplementationTypes(OrderProcessingWorkflowImpl.class);
        orderProcessingWorker.registerActivitiesImplementations(new OrderProcessingActivityImpl());
        
        // Additional workers can be registered here for different task queues
        // Example:
        // Worker paymentWorker = workerFactory.newWorker(paymentProcessingQueue);
        // paymentWorker.registerWorkflowImplementationTypes(PaymentWorkflowImpl.class);
        // paymentWorker.registerActivitiesImplementations(new PaymentActivityImpl());
    }
}
EOF
    fi
fi

# Rename docker-compose.yml.tmpl to compose.yaml.tmpl for consistency
if [ -f "$TEMPLATE_DIR/docker-compose.yml.tmpl" ] && [ ! -f "$TEMPLATE_DIR/compose.yaml.tmpl" ]; then
    echo "Renaming docker-compose.yml.tmpl to compose.yaml.tmpl"
    mv "$TEMPLATE_DIR/docker-compose.yml.tmpl" "$TEMPLATE_DIR/compose.yaml.tmpl"
fi

# Create .github/workflows directory and move github workflow file if it exists
mkdir -p "$TEMPLATE_DIR/.github/workflows"
if [ -f "$TEMPLATE_DIR/github-workflow.yml.tmpl" ]; then
    echo "Moving github-workflow.yml.tmpl to .github/workflows directory"
    mv "$TEMPLATE_DIR/github-workflow.yml.tmpl" "$TEMPLATE_DIR/.github/workflows/ci-cd.yml.tmpl"
fi

# Create charts directory with basic structure for Helm charts
mkdir -p "$TEMPLATE_DIR/charts/{{artifactId}}/templates"
if [ ! -f "$TEMPLATE_DIR/charts/{{artifactId}}/Chart.yaml.tmpl" ]; then
    echo "Creating Helm Chart.yaml template"
    cat > "$TEMPLATE_DIR/charts/{{artifactId}}/Chart.yaml.tmpl" << EOF
apiVersion: v2
name: {{artifactId}}
description: A Helm chart for {{name}} application
type: application
version: 0.1.0
appVersion: "1.0.0"
EOF
fi

if [ ! -f "$TEMPLATE_DIR/charts/{{artifactId}}/values.yaml.tmpl" ]; then
    echo "Creating Helm values.yaml template"
    cat > "$TEMPLATE_DIR/charts/{{artifactId}}/values.yaml.tmpl" << EOF
# Default values for {{artifactId}}.
# This is a YAML-formatted file.

replicaCount: 1

image:
  repository: {{artifactId}}
  pullPolicy: IfNotPresent
  tag: "latest"

nameOverride: ""
fullnameOverride: ""

serviceAccount:
  create: true
  name: ""

service:
  type: ClusterIP
  port: 8080

ingress:
  enabled: false
  annotations: {}
  hosts:
    - host: chart-example.local
      paths: []
  tls: []

resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 500m
    memory: 512Mi

autoscaling:
  enabled: false
  minReplicas: 1
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80

nodeSelector: {}

tolerations: []

affinity: {}

# Application specific configuration
env:
  SPRING_PROFILES_ACTIVE: prod
  TEMPORAL_SERVICE_ADDRESS: temporal.default.svc.cluster.local:7233
  TEMPORAL_NAMESPACE: default
  AWS_REGION: us-west-2

# Database configuration
database:
  host: postgres.default.svc.cluster.local
  port: 5432
  name: {{artifactId}}
  username: postgres
  existingSecret: {{artifactId}}-db-credentials
  existingSecretKey: password
EOF
fi

# Create deployment.yaml template
mkdir -p "$TEMPLATE_DIR/charts/{{artifactId}}/templates"
if [ ! -f "$TEMPLATE_DIR/charts/{{artifactId}}/templates/deployment.yaml.tmpl" ]; then
    echo "Creating Helm deployment.yaml template"
    cat > "$TEMPLATE_DIR/charts/{{artifactId}}/templates/deployment.yaml.tmpl" << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "{{artifactId}}.fullname" . }}
  labels:
    {{- include "{{artifactId}}.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "{{artifactId}}.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "{{artifactId}}.selectorLabels" . | nindent 8 }}
    spec:
      serviceAccountName: {{ include "{{artifactId}}.serviceAccountName" . }}
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
          env:
            - name: SPRING_PROFILES_ACTIVE
              value: {{ .Values.env.SPRING_PROFILES_ACTIVE }}
            - name: SPRING_DATASOURCE_URL
              value: jdbc:postgresql://{{ .Values.database.host }}:{{ .Values.database.port }}/{{ .Values.database.name }}
            - name: SPRING_DATASOURCE_USERNAME
              value: {{ .Values.database.username }}
            - name: SPRING_DATASOURCE_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: {{ .Values.database.existingSecret }}
                  key: {{ .Values.database.existingSecretKey }}
            - name: TEMPORAL_SERVICE_ADDRESS
              value: {{ .Values.env.TEMPORAL_SERVICE_ADDRESS }}
            - name: TEMPORAL_NAMESPACE
              value: {{ .Values.env.TEMPORAL_NAMESPACE }}
            - name: AWS_REGION
              value: {{ .Values.env.AWS_REGION }}
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
EOF
fi

# Create a basic controller example if it doesn't exist
if [ ! -f "$TEMPLATE_DIR/controller/HealthController.java.tmpl" ]; then
    echo "Creating HealthController.java.tmpl"
    cat > "$TEMPLATE_DIR/controller/HealthController.java.tmpl" << EOF
package {{package}}.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.Map;

/**
 * Controller for health checks and service status.
 */
@RestController
@RequestMapping("/health")
public class HealthController {

    /**
     * Returns the health status of the application.
     *
     * @return A response with health status information
     */
    @GetMapping
    public ResponseEntity<Map<String, Object>> healthCheck() {
        Map<String, Object> healthStatus = new HashMap<>();
        healthStatus.put("status", "UP");
        healthStatus.put("timestamp", LocalDateTime.now().toString());
        healthStatus.put("service", "{{name}}");
        
        return ResponseEntity.ok(healthStatus);
    }
}
EOF
fi

# Create a basic flyway migration script example
mkdir -p "$TEMPLATE_DIR/src/main/resources/db/migration"
if [ ! -f "$TEMPLATE_DIR/src/main/resources/db/migration/V1__init_schema.sql.tmpl" ]; then
    echo "Creating initial Flyway migration script"
    cat > "$TEMPLATE_DIR/src/main/resources/db/migration/V1__init_schema.sql.tmpl" << EOF
-- Initial database schema for {{artifactId}}

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(100) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    last_login_at TIMESTAMP,
    active BOOLEAN NOT NULL DEFAULT TRUE
);

-- Create indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- Example entity table (orders)
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    order_number VARCHAR(50) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    
    CONSTRAINT chk_status CHECK (status IN ('PENDING', 'PROCESSING', 'COMPLETED', 'CANCELLED', 'FAILED'))
);

-- Create indexes
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
EOF
fi

# Create OpenAPI spec if it doesn't exist
mkdir -p "$TEMPLATE_DIR/src/main/swagger"
if [ ! -f "$TEMPLATE_DIR/src/main/swagger/api-spec.yaml.tmpl" ]; then
    echo "Creating OpenAPI specification"
    cat > "$TEMPLATE_DIR/src/main/swagger/api-spec.yaml.tmpl" << EOF
openapi: 3.0.3
info:
  title: {{name}} API
  description: API documentation for {{name}}
  version: 1.0.0
  contact:
    name: Development Team
  license:
    name: Private
servers:
  - url: /api
    description: API base path
tags:
  - name: Health
    description: Health check endpoints
  - name: Auth
    description: Authentication endpoints
  - name: Orders
    description: Order management endpoints
paths:
  /health:
    get:
      tags:
        - Health
      summary: Health check
      description: Returns the health status of the service
      operationId: healthCheck
      responses:
        '200':
          description: Successful operation
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                  timestamp:
                    type: string
                  service:
                    type: string
  /auth/login:
    post:
      tags:
        - Auth
      summary: User login
      description: Authenticates a user and returns a JWT token
      operationId: login
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required:
                - username
                - password
              properties:
                username:
                  type: string
                password:
                  type: string
      responses:
        '200':
          description: Successfully authenticated
          content:
            application/json:
              schema:
                type: object
                properties:
                  token:
                    type: string
                  expiresIn:
                    type: integer
        '401':
          description: Authentication failed
  /orders:
    get:
      tags:
        - Orders
      summary: Get all orders
      description: Returns a list of all orders for the authenticated user
      operationId: getOrders
      security:
        - bearerAuth: []
      responses:
        '200':
          description: Successful operation
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Order'
        '401':
          description: Unauthorized
    post:
      tags:
        - Orders
      summary: Create a new order
      description: Creates a new order for the authenticated user
      operationId: createOrder
      security:
        - bearerAuth: []
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/OrderRequest'
      responses:
        '201':
          description: Order created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Order'
        '400':
          description: Invalid request
        '401':
          description: Unauthorized
components:
  schemas:
    Order:
      type: object
      properties:
        id:
          type: string
          format: uuid
        orderNumber:
          type: string
        status:
          type: string
          enum:
            - PENDING
            - PROCESSING
            - COMPLETED
            - CANCELLED
            - FAILED
        totalAmount:
          type: number
          format: double
        createdAt:
          type: string
          format: date-time
        updatedAt:
          type: string
          format: date-time
    OrderRequest:
      type: object
      required:
        - items
      properties:
        items:
          type: array
          items:
            type: object
            properties:
              productId:
                type: string
                format: uuid
              quantity:
                type: integer
                minimum: 1
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
EOF
fi

echo "Template organization complete!" 
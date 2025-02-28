package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/springwell/cli/pkg/config"
	"github.com/springwell/cli/pkg/util"
	"github.com/urfave/cli/v2"
)

// InteractiveCommand returns the command to run the CLI in interactive mode
func InteractiveCommand() *cli.Command {
	return &cli.Command{
		Name:    "interactive",
		Aliases: []string{"i"},
		Usage:   "Run the CLI in interactive mode",
		Action:  runInteractiveMode,
	}
}

// runInteractiveMode runs the CLI in interactive mode
func runInteractiveMode(c *cli.Context) error {
	reader := bufio.NewReader(os.Stdin)

	// Print welcome message
	printWelcome()

	for {
		// Display main menu
		printMainMenu()

		// Get user selection
		fmt.Print("Enter your choice: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		// Trim newline characters
		input = strings.TrimSpace(input)

		// Process selection
		if err := processMainMenuSelection(input, reader); err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		// Exit if user selected exit option
		if input == "0" {
			break
		}
	}

	return nil
}

// printWelcome prints the welcome message
func printWelcome() {
	util.PrintSuccess("\n=================================")
	util.PrintSuccess("  Welcome to SpringWell CLI")
	util.PrintSuccess("  Your Spring Boot Companion")
	util.PrintSuccess("=================================\n")
}

// printMainMenu prints the main menu options
func printMainMenu() {
	fmt.Println("\nMain Menu:")
	fmt.Println("1. Generate a brand new project")
	fmt.Println("2. Generate components (entity, workflow, etc.)")
	fmt.Println("3. Run development server")
	fmt.Println("4. Build project")
	fmt.Println("5. Run tests")
	fmt.Println("6. Check project health")
	fmt.Println("0. Exit")
}

// processMainMenuSelection processes the user's selection from the main menu
func processMainMenuSelection(input string, reader *bufio.Reader) error {
	switch input {
	case "0":
		util.PrintInfo("Exiting SpringWell CLI. Goodbye!")
		return nil
	case "1":
		return handleGenerateProject(reader)
	case "2":
		return handleGenerateComponents(reader)
	case "3":
		return handleRunDev(reader)
	case "4":
		return handleBuild()
	case "5":
		return handleRunTests(reader)
	case "6":
		return handleProjectHealth()
	default:
		fmt.Println("Invalid selection. Please try again.")
		return nil
	}
}

// handleGenerateProject handles the "Generate a brand new project" option
func handleGenerateProject(reader *bufio.Reader) error {
	fmt.Println("\n=== Generate New Project ===")

	// Get project name
	fmt.Print("Enter project name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("Project name is required")
	}

	// Get package name (optional)
	fmt.Print("Enter package name (leave blank for default): ")
	packageName, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	packageName = strings.TrimSpace(packageName)

	if packageName == "" {
		packageName = "com." + util.ToPackageName(name)
	}

	// Select template
	fmt.Println("\nSelect project template:")
	fmt.Println("1. Basic Spring Boot")
	fmt.Println("2. AWS + Temporal + Auth0")
	fmt.Print("Enter your choice (default: 1): ")

	templateChoice, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	templateChoice = strings.TrimSpace(templateChoice)

	template := "basic"
	if templateChoice == "2" {
		template = "aws-temporal-auth0"
	}

	// Select database
	fmt.Println("\nSelect database:")
	fmt.Println("1. PostgreSQL")
	fmt.Println("2. MySQL")
	fmt.Println("3. H2 (in-memory)")
	fmt.Print("Enter your choice (default: 1): ")

	dbChoice, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	dbChoice = strings.TrimSpace(dbChoice)

	db := "postgres"
	switch dbChoice {
	case "2":
		db = "mysql"
	case "3":
		db = "h2"
	}

	// Create project directory
	projectDir := filepath.Join(".", name)
	if err := util.CreateDirectory(projectDir); err != nil {
		return err
	}

	// Set up configuration
	cfg := config.GetDefaultConfig()
	cfg.Project.Package = packageName

	// Save configuration
	if err := config.SaveConfig(cfg, projectDir); err != nil {
		return err
	}

	// Create the project
	util.PrintInfo("\nCreating project %s with template %s...", name, template)

	if template == "aws-temporal-auth0" {
		return createAwsTemporalAuth0Project(name, packageName, projectDir, db)
	}

	return createSpringBootProject(name, packageName, projectDir, db, "jwt", "swagger,actuator")
}

// handleGenerateComponents handles the "Generate components" option
func handleGenerateComponents(reader *bufio.Reader) error {
	fmt.Println("\n=== Generate Components ===")

	// Check if it's a Spring Boot project
	if !util.IsSpringBootProject(".") {
		return fmt.Errorf("Current directory is not a Spring Boot project")
	}

	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return err
	}

	// Display component options
	fmt.Println("\nSelect component to generate:")
	fmt.Println("1. Entity")
	fmt.Println("2. Controller")
	fmt.Println("3. Service")
	fmt.Println("4. Repository")
	fmt.Println("5. DTO")
	fmt.Println("6. Temporal Workflow")
	fmt.Println("7. Temporal Activity")
	fmt.Println("0. Back to main menu")

	fmt.Print("Enter your choice: ")
	componentChoice, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	componentChoice = strings.TrimSpace(componentChoice)

	if componentChoice == "0" {
		return nil
	}

	// Get component name
	fmt.Print("Enter component name: ")
	componentName, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	componentName = strings.TrimSpace(componentName)

	if componentName == "" {
		return fmt.Errorf("Component name is required")
	}

	// Process component generation based on choice
	switch componentChoice {
	case "1":
		return generateEntity(componentName, cfg.Project.Package)
	case "2":
		return generateController(componentName, cfg.Project.Package)
	case "3":
		return generateService(componentName, cfg.Project.Package)
	case "4":
		return generateRepository(componentName, cfg.Project.Package)
	case "5":
		return generateDTO(componentName, cfg.Project.Package)
	case "6":
		return generateTemporalWorkflow(componentName, cfg.Project.Package)
	case "7":
		return generateTemporalActivity(componentName, cfg.Project.Package)
	default:
		return fmt.Errorf("Invalid component choice")
	}
}

// handleRunDev handles the "Run development server" option
func handleRunDev(reader *bufio.Reader) error {
	fmt.Println("\n=== Run Development Server ===")

	// Check if it's a Spring Boot project
	if !util.IsSpringBootProject(".") {
		return fmt.Errorf("Current directory is not a Spring Boot project")
	}

	// Get port (optional)
	fmt.Print("Enter port (default: 8080): ")
	portStr, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	portStr = strings.TrimSpace(portStr)

	port := "8080"
	if portStr != "" {
		port = portStr
	}

	// Get profile (optional)
	fmt.Print("Enter profile (default: dev): ")
	profile, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	profile = strings.TrimSpace(profile)

	if profile == "" {
		profile = "dev"
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
			"-Dspring-boot.run.profiles="+profile,
			"-Dserver.port="+port)
	} else {
		cmd = exec.Command("./gradlew", "bootRun",
			"-Dspring.profiles.active="+profile,
			"-Dserver.port="+port)
	}

	// Set up command output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	util.PrintInfo("Starting application in development mode... (Ctrl+C to stop)")
	return cmd.Run()
}

// handleBuild handles the "Build project" option
func handleBuild() error {
	fmt.Println("\n=== Build Project ===")

	// Check if it's a Spring Boot project
	if !util.IsSpringBootProject(".") {
		return fmt.Errorf("Current directory is not a Spring Boot project")
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
}

// handleRunTests handles the "Run tests" option
func handleRunTests(reader *bufio.Reader) error {
	fmt.Println("\n=== Run Tests ===")

	// Check if it's a Spring Boot project
	if !util.IsSpringBootProject(".") {
		return fmt.Errorf("Current directory is not a Spring Boot project")
	}

	// Ask for specific test (optional)
	fmt.Print("Enter specific test to run (leave blank for all tests): ")
	testName, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	testName = strings.TrimSpace(testName)

	// Determine if it's Maven or Gradle
	isMaven := true
	if _, err := os.Stat("build.gradle"); err == nil {
		isMaven = false
	}

	// Build the command
	var cmd *exec.Cmd
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
}

// handleProjectHealth handles the "Check project health" option
func handleProjectHealth() error {
	fmt.Println("\n=== Project Health Check ===")

	// Check if it's a Spring Boot project
	if !util.IsSpringBootProject(".") {
		return fmt.Errorf("Current directory is not a Spring Boot project")
	}

	// TODO: Implement more comprehensive health checks

	// Check for common files and directories
	essentialFiles := []string{
		"pom.xml", "build.gradle",
		"src/main/java",
		"src/main/resources",
		"src/test",
	}

	missingFiles := []string{}
	for _, file := range essentialFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			missingFiles = append(missingFiles, file)
		}
	}

	if len(missingFiles) > 0 {
		util.PrintWarning("Missing files/directories:")
		for _, file := range missingFiles {
			fmt.Printf("  - %s\n", file)
		}
	}

	// Check application properties
	propertiesFiles := []string{
		"src/main/resources/application.properties",
		"src/main/resources/application.yml",
	}

	propertiesFound := false
	for _, file := range propertiesFiles {
		if _, err := os.Stat(file); err == nil {
			propertiesFound = true
			break
		}
	}

	if !propertiesFound {
		util.PrintWarning("No application properties file found")
	}

	util.PrintSuccess("Project health check completed")
	return nil
}

// generateEntity creates a new entity class
func generateEntity(name string, packageName string) error {
	// Capitalize first letter for class name
	className := util.ToCamelCase(name)

	// Create entity directory if it doesn't exist
	entityDir := filepath.Join("src/main/java", strings.ReplaceAll(packageName, ".", "/"), "domain/entity")
	if err := util.CreateDirectory(entityDir); err != nil {
		return err
	}

	// Entity template
	entityTemplate := `package ${package}.domain.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import java.time.LocalDateTime;

@Entity
@Table(name = "${tableName}")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class ${className} {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    private String name;
    
    private String description;
    
    @CreationTimestamp
    @Column(name = "created_at", updatable = false)
    private LocalDateTime createdAt;
    
    @UpdateTimestamp
    @Column(name = "updated_at")
    private LocalDateTime updatedAt;
}
`

	// Replace placeholders
	entityTemplate = strings.ReplaceAll(entityTemplate, "${package}", packageName)
	entityTemplate = strings.ReplaceAll(entityTemplate, "${className}", className)
	entityTemplate = strings.ReplaceAll(entityTemplate, "${tableName}", util.ToSnakeCase(name)+"s")

	// Write to file
	entityFile := filepath.Join(entityDir, className+".java")
	if err := util.WriteFile(entityFile, entityTemplate); err != nil {
		return err
	}

	util.PrintSuccess("Generated entity: %s", entityFile)
	return nil
}

// generateController creates a new controller class
func generateController(name string, packageName string) error {
	// Capitalize first letter for class name
	className := util.ToCamelCase(name)

	// Create controller directory if it doesn't exist
	controllerDir := filepath.Join("src/main/java", strings.ReplaceAll(packageName, ".", "/"), "controller")
	if err := util.CreateDirectory(controllerDir); err != nil {
		return err
	}

	// Controller template
	controllerTemplate := `package ${package}.controller;

import ${package}.domain.entity.${className};
import ${package}.service.${className}Service;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/${resourceName}")
@RequiredArgsConstructor
public class ${className}Controller {

    private final ${className}Service ${variableName}Service;
    
    @GetMapping
    public ResponseEntity<List<${className}>> getAll() {
        return ResponseEntity.ok(${variableName}Service.findAll());
    }
    
    @GetMapping("/{id}")
    public ResponseEntity<${className}> getById(@PathVariable Long id) {
        return ${variableName}Service.findById(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @PostMapping
    public ResponseEntity<${className}> create(@RequestBody ${className} ${variableName}) {
        return ResponseEntity.ok(${variableName}Service.save(${variableName}));
    }
    
    @PutMapping("/{id}")
    public ResponseEntity<${className}> update(@PathVariable Long id, @RequestBody ${className} ${variableName}) {
        ${variableName}.setId(id);
        return ResponseEntity.ok(${variableName}Service.save(${variableName}));
    }
    
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> delete(@PathVariable Long id) {
        ${variableName}Service.deleteById(id);
        return ResponseEntity.noContent().build();
    }
}
`

	// Replace placeholders
	controllerTemplate = strings.ReplaceAll(controllerTemplate, "${package}", packageName)
	controllerTemplate = strings.ReplaceAll(controllerTemplate, "${className}", className)
	controllerTemplate = strings.ReplaceAll(controllerTemplate, "${variableName}", util.ToCamelCaseFirstLower(className))
	controllerTemplate = strings.ReplaceAll(controllerTemplate, "${resourceName}", util.ToKebabCase(name)+"s")

	// Write to file
	controllerFile := filepath.Join(controllerDir, className+"Controller.java")
	if err := util.WriteFile(controllerFile, controllerTemplate); err != nil {
		return err
	}

	util.PrintSuccess("Generated controller: %s", controllerFile)
	return nil
}

// generateService creates a new service class
func generateService(name string, packageName string) error {
	// Capitalize first letter for class name
	className := util.ToCamelCase(name)

	// Create service directory if it doesn't exist
	serviceDir := filepath.Join("src/main/java", strings.ReplaceAll(packageName, ".", "/"), "service")
	if err := util.CreateDirectory(serviceDir); err != nil {
		return err
	}

	// Service interface template
	serviceInterfaceTemplate := `package ${package}.service;

import ${package}.domain.entity.${className};

import java.util.List;
import java.util.Optional;

public interface ${className}Service {

    List<${className}> findAll();
    
    Optional<${className}> findById(Long id);
    
    ${className} save(${className} ${variableName});
    
    void deleteById(Long id);
}
`

	// Service implementation template
	serviceImplTemplate := `package ${package}.service.impl;

import ${package}.domain.entity.${className};
import ${package}.repository.${className}Repository;
import ${package}.service.${className}Service;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
@Transactional
public class ${className}ServiceImpl implements ${className}Service {

    private final ${className}Repository ${variableName}Repository;
    
    @Override
    @Transactional(readOnly = true)
    public List<${className}> findAll() {
        return ${variableName}Repository.findAll();
    }
    
    @Override
    @Transactional(readOnly = true)
    public Optional<${className}> findById(Long id) {
        return ${variableName}Repository.findById(id);
    }
    
    @Override
    public ${className} save(${className} ${variableName}) {
        return ${variableName}Repository.save(${variableName});
    }
    
    @Override
    public void deleteById(Long id) {
        ${variableName}Repository.deleteById(id);
    }
}
`

	// Replace placeholders
	serviceInterfaceTemplate = strings.ReplaceAll(serviceInterfaceTemplate, "${package}", packageName)
	serviceInterfaceTemplate = strings.ReplaceAll(serviceInterfaceTemplate, "${className}", className)
	serviceInterfaceTemplate = strings.ReplaceAll(serviceInterfaceTemplate, "${variableName}", util.ToCamelCaseFirstLower(className))

	serviceImplTemplate = strings.ReplaceAll(serviceImplTemplate, "${package}", packageName)
	serviceImplTemplate = strings.ReplaceAll(serviceImplTemplate, "${className}", className)
	serviceImplTemplate = strings.ReplaceAll(serviceImplTemplate, "${variableName}", util.ToCamelCaseFirstLower(className))

	// Create impl directory if it doesn't exist
	serviceImplDir := filepath.Join(serviceDir, "impl")
	if err := util.CreateDirectory(serviceImplDir); err != nil {
		return err
	}

	// Write to files
	serviceInterfaceFile := filepath.Join(serviceDir, className+"Service.java")
	if err := util.WriteFile(serviceInterfaceFile, serviceInterfaceTemplate); err != nil {
		return err
	}

	serviceImplFile := filepath.Join(serviceImplDir, className+"ServiceImpl.java")
	if err := util.WriteFile(serviceImplFile, serviceImplTemplate); err != nil {
		return err
	}

	util.PrintSuccess("Generated service interface: %s", serviceInterfaceFile)
	util.PrintSuccess("Generated service implementation: %s", serviceImplFile)
	return nil
}

// generateRepository creates a new repository interface
func generateRepository(name string, packageName string) error {
	// Capitalize first letter for class name
	className := util.ToCamelCase(name)

	// Create repository directory if it doesn't exist
	repositoryDir := filepath.Join("src/main/java", strings.ReplaceAll(packageName, ".", "/"), "repository")
	if err := util.CreateDirectory(repositoryDir); err != nil {
		return err
	}

	// Repository template
	repositoryTemplate := `package ${package}.repository;

import ${package}.domain.entity.${className};
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface ${className}Repository extends JpaRepository<${className}, Long> {
    // Add custom query methods here
}
`

	// Replace placeholders
	repositoryTemplate = strings.ReplaceAll(repositoryTemplate, "${package}", packageName)
	repositoryTemplate = strings.ReplaceAll(repositoryTemplate, "${className}", className)

	// Write to file
	repositoryFile := filepath.Join(repositoryDir, className+"Repository.java")
	if err := util.WriteFile(repositoryFile, repositoryTemplate); err != nil {
		return err
	}

	util.PrintSuccess("Generated repository: %s", repositoryFile)
	return nil
}

// generateDTO creates a new DTO class
func generateDTO(name string, packageName string) error {
	// Capitalize first letter for class name
	className := util.ToCamelCase(name)

	// Create DTO directory if it doesn't exist
	dtoDir := filepath.Join("src/main/java", strings.ReplaceAll(packageName, ".", "/"), "domain/dto")
	if err := util.CreateDirectory(dtoDir); err != nil {
		return err
	}

	// DTO template
	dtoTemplate := `package ${package}.domain.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ${className}DTO {
    
    private Long id;
    
    private String name;
    
    private String description;
    
    private LocalDateTime createdAt;
    
    private LocalDateTime updatedAt;
}
`

	// Replace placeholders
	dtoTemplate = strings.ReplaceAll(dtoTemplate, "${package}", packageName)
	dtoTemplate = strings.ReplaceAll(dtoTemplate, "${className}", className)

	// Write to file
	dtoFile := filepath.Join(dtoDir, className+"DTO.java")
	if err := util.WriteFile(dtoFile, dtoTemplate); err != nil {
		return err
	}

	util.PrintSuccess("Generated DTO: %s", dtoFile)
	return nil
}

// generateTemporalWorkflow creates new Temporal workflow files
func generateTemporalWorkflow(name string, packageName string) error {
	// Capitalize first letter for class name
	className := util.ToCamelCase(name)

	// Create workflow directories if they don't exist
	workflowDir := filepath.Join("src/main/java", strings.ReplaceAll(packageName, ".", "/"), "temporal/workflow")
	if err := util.CreateDirectory(workflowDir); err != nil {
		return err
	}

	workflowImplDir := filepath.Join("src/main/java", strings.ReplaceAll(packageName, ".", "/"), "temporal/workflow/impl")
	if err := util.CreateDirectory(workflowImplDir); err != nil {
		return err
	}

	// Workflow interface template
	workflowTemplate := `package ${package}.temporal.workflow;

import io.temporal.workflow.WorkflowInterface;
import io.temporal.workflow.WorkflowMethod;

/**
 * ${className} workflow interface for Temporal.
 */
@WorkflowInterface
public interface ${className}Workflow {

    /**
     * Main workflow method.
     * 
     * @param input The input for the workflow
     * @return The result of the workflow execution
     */
    @WorkflowMethod
    String execute(String input);
}
`

	// Workflow implementation template
	workflowImplTemplate := `package ${package}.temporal.workflow.impl;

import ${package}.temporal.workflow.${className}Workflow;
import ${package}.temporal.activity.${className}Activity;

import io.temporal.activity.ActivityOptions;
import io.temporal.common.RetryOptions;
import io.temporal.workflow.Workflow;

import java.time.Duration;

/**
 * Implementation of the ${className}Workflow interface.
 */
public class ${className}WorkflowImpl implements ${className}Workflow {

    private final ActivityOptions activityOptions = ActivityOptions.newBuilder()
            .setScheduleToCloseTimeout(Duration.ofSeconds(10))
            .setRetryOptions(RetryOptions.newBuilder()
                    .setMaximumAttempts(3)
                    .build())
            .build();

    private final ${className}Activity activity = Workflow.newActivityStub(
            ${className}Activity.class, activityOptions);

    @Override
    public String execute(String input) {
        // Log the start of the workflow
        Workflow.getLogger(${className}WorkflowImpl.class)
                .info("${className} workflow execution started with input: " + input);

        // Execute an activity as part of this workflow
        String result = activity.execute(input);

        // Log the completion of the workflow
        Workflow.getLogger(${className}WorkflowImpl.class)
                .info("${className} workflow execution completed with result: " + result);

        return result;
    }
}
`

	// Replace placeholders
	workflowTemplate = strings.ReplaceAll(workflowTemplate, "${package}", packageName)
	workflowTemplate = strings.ReplaceAll(workflowTemplate, "${className}", className)

	workflowImplTemplate = strings.ReplaceAll(workflowImplTemplate, "${package}", packageName)
	workflowImplTemplate = strings.ReplaceAll(workflowImplTemplate, "${className}", className)

	// Write to files
	workflowFile := filepath.Join(workflowDir, className+"Workflow.java")
	if err := util.WriteFile(workflowFile, workflowTemplate); err != nil {
		return err
	}

	workflowImplFile := filepath.Join(workflowImplDir, className+"WorkflowImpl.java")
	if err := util.WriteFile(workflowImplFile, workflowImplTemplate); err != nil {
		return err
	}

	util.PrintSuccess("Generated workflow interface: %s", workflowFile)
	util.PrintSuccess("Generated workflow implementation: %s", workflowImplFile)

	// Also generate the corresponding activity
	return generateTemporalActivity(name, packageName)
}

// generateTemporalActivity creates new Temporal activity files
func generateTemporalActivity(name string, packageName string) error {
	// Capitalize first letter for class name
	className := util.ToCamelCase(name)

	// Create activity directories if they don't exist
	activityDir := filepath.Join("src/main/java", strings.ReplaceAll(packageName, ".", "/"), "temporal/activity")
	if err := util.CreateDirectory(activityDir); err != nil {
		return err
	}

	activityImplDir := filepath.Join("src/main/java", strings.ReplaceAll(packageName, ".", "/"), "temporal/activity/impl")
	if err := util.CreateDirectory(activityImplDir); err != nil {
		return err
	}

	// Activity interface template
	activityTemplate := `package ${package}.temporal.activity;

import io.temporal.activity.ActivityInterface;
import io.temporal.activity.ActivityMethod;

/**
 * ${className} activity interface for Temporal.
 */
@ActivityInterface
public interface ${className}Activity {

    /**
     * Main activity method.
     * 
     * @param input The input for the activity
     * @return The result of the activity execution
     */
    @ActivityMethod
    String execute(String input);
}
`

	// Activity implementation template
	activityImplTemplate := `package ${package}.temporal.activity.impl;

import ${package}.temporal.activity.${className}Activity;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

/**
 * Implementation of the ${className}Activity interface.
 */
@Component
public class ${className}ActivityImpl implements ${className}Activity {

    private static final Logger logger = LoggerFactory.getLogger(${className}ActivityImpl.class);

    @Override
    public String execute(String input) {
        logger.info("Executing ${className} activity with input: {}", input);
        
        // Implement your business logic here
        String result = processInput(input);
        
        logger.info("${className} activity execution completed with result: {}", result);
        return result;
    }
    
    private String processInput(String input) {
        // This is where you would implement your business logic
        // For now, we'll just append a message
        return input + " - processed by ${className}Activity";
    }
}
`

	// Replace placeholders
	activityTemplate = strings.ReplaceAll(activityTemplate, "${package}", packageName)
	activityTemplate = strings.ReplaceAll(activityTemplate, "${className}", className)

	activityImplTemplate = strings.ReplaceAll(activityImplTemplate, "${package}", packageName)
	activityImplTemplate = strings.ReplaceAll(activityImplTemplate, "${className}", className)

	// Write to files
	activityFile := filepath.Join(activityDir, className+"Activity.java")
	if err := util.WriteFile(activityFile, activityTemplate); err != nil {
		return err
	}

	activityImplFile := filepath.Join(activityImplDir, className+"ActivityImpl.java")
	if err := util.WriteFile(activityImplFile, activityImplTemplate); err != nil {
		return err
	}

	util.PrintSuccess("Generated activity interface: %s", activityFile)
	util.PrintSuccess("Generated activity implementation: %s", activityImplFile)

	return nil
}

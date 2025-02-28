package util

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// Colors
var (
	InfoColor    = color.New(color.FgCyan)
	SuccessColor = color.New(color.FgGreen)
	ErrorColor   = color.New(color.FgRed)
	WarnColor    = color.New(color.FgYellow)
	BoldColor    = color.New(color.Bold)
)

// PrintInfo prints an information message
func PrintInfo(format string, a ...interface{}) {
	InfoColor.Printf(format+"\n", a...)
}

// PrintSuccess prints a success message
func PrintSuccess(format string, a ...interface{}) {
	SuccessColor.Printf("✓ "+format+"\n", a...)
}

// PrintError prints an error message
func PrintError(format string, a ...interface{}) {
	ErrorColor.Printf("✗ "+format+"\n", a...)
}

// PrintWarning prints a warning message
func PrintWarning(format string, a ...interface{}) {
	WarnColor.Printf("! "+format+"\n", a...)
}

// PrintBold prints a bold message
func PrintBold(format string, a ...interface{}) {
	BoldColor.Printf(format+"\n", a...)
}

// IsSpringBootProject checks if the current directory is a Spring Boot project
func IsSpringBootProject(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "pom.xml"))
	if err == nil {
		return true
	}

	_, err = os.Stat(filepath.Join(dir, "build.gradle"))
	if err == nil {
		return true
	}

	return false
}

// CreateDirectory creates a directory if it doesn't exist
func CreateDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// WriteFile writes content to a file, creating directories as needed
func WriteFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := CreateDirectory(dir); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// ToJavaClassName converts a string to a Java class name (PascalCase)
func ToJavaClassName(name string) string {
	words := splitByCase(name)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[0:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}

// ToJavaVariableName converts a string to a Java variable name (camelCase)
func ToJavaVariableName(name string) string {
	className := ToJavaClassName(name)
	if len(className) > 0 {
		return strings.ToLower(className[0:1]) + className[1:]
	}
	return ""
}

// ToDatabaseTableName converts a string to a database table name (snake_case)
func ToDatabaseTableName(name string) string {
	words := splitByCase(name)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "_")
}

// ToColumnName converts a field name to a database column name (snake_case)
func ToColumnName(name string) string {
	return ToDatabaseTableName(name)
}

// ToPackageName converts a string to a package name (lowercase with dots)
func ToPackageName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", ".")
	name = strings.ReplaceAll(name, "_", ".")
	return name
}

// splitByCase splits a string by case (camelCase, PascalCase, snake_case, kebab-case)
func splitByCase(s string) []string {
	// Replace kebab-case and snake_case with spaces
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	// Insert space before capital letters (for camelCase and PascalCase)
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	s = re.ReplaceAllString(s, "$1 $2")

	// Split by spaces
	words := strings.Fields(s)
	return words
}

// ParseFieldDefinitions parses field definitions from a string
// Format: "name:type[:modifier]"
func ParseFieldDefinitions(fields string) ([]map[string]string, error) {
	if fields == "" {
		return []map[string]string{}, nil
	}

	var result []map[string]string
	fieldsList := strings.Split(fields, " ")

	for _, field := range fieldsList {
		parts := strings.Split(field, ":")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid field format: %s, expected name:type[:modifier]", field)
		}

		fieldMap := map[string]string{
			"name": parts[0],
			"type": parts[1],
		}

		// Add column name
		fieldMap["columnName"] = ToColumnName(parts[0])

		// Add nullable modifier if specified
		if len(parts) > 2 && parts[2] == "nullable" {
			fieldMap["nullable"] = "true"
		} else {
			fieldMap["nullable"] = "false"
		}

		result = append(result, fieldMap)
	}

	return result, nil
}

// ParseRelationships parses relationship definitions from a string
// Format: "type:field:entity"
func ParseRelationships(relations string) ([]map[string]string, error) {
	if relations == "" {
		return []map[string]string{}, nil
	}

	var result []map[string]string
	relationsList := strings.Split(relations, " ")

	for _, relation := range relationsList {
		parts := strings.Split(relation, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid relationship format: %s, expected type:field:entity", relation)
		}

		relationType := parts[0]
		field := parts[1]
		entity := parts[2]

		// Validate relationship type
		validTypes := map[string]bool{
			"oneToOne":   true,
			"oneToMany":  true,
			"manyToOne":  true,
			"manyToMany": true,
		}

		if !validTypes[relationType] {
			return nil, fmt.Errorf("invalid relationship type: %s, expected oneToOne, oneToMany, manyToOne, or manyToMany", relationType)
		}

		relationMap := map[string]string{
			"type":   relationType,
			"field":  field,
			"entity": entity,
		}

		result = append(result, relationMap)
	}

	return result, nil
}

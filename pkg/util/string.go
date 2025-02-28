package util

import (
	"regexp"
	"strings"
)

// ToCamelCase converts a string to CamelCase (PascalCase)
func ToCamelCase(s string) string {
	// Replace non-alphanumeric characters with spaces
	s = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(s, " ")

	// Title case each word
	s = strings.Title(strings.ToLower(s))

	// Remove spaces
	s = strings.ReplaceAll(s, " ", "")

	return s
}

// ToCamelCaseFirstLower converts a string to camelCase (first letter lowercase)
func ToCamelCaseFirstLower(s string) string {
	s = ToCamelCase(s)
	if len(s) == 0 {
		return s
	}

	return strings.ToLower(s[:1]) + s[1:]
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	// Replace non-alphanumeric characters with spaces
	s = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(s, " ")

	// Replace spaces with underscores and convert to lowercase
	s = strings.ToLower(strings.ReplaceAll(s, " ", "_"))

	return s
}

// ToKebabCase converts a string to kebab-case
func ToKebabCase(s string) string {
	// Replace non-alphanumeric characters with spaces
	s = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(s, " ")

	// Replace spaces with hyphens and convert to lowercase
	s = strings.ToLower(strings.ReplaceAll(s, " ", "-"))

	return s
}

package rxtspot

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	// MaxNameLength is the maximum allowed length for resource names
	MaxNameLength = 63
	// DefaultRequestTimeout is the default timeout for API requests
	DefaultRequestTimeout = 30 * time.Second
)

var (
	dns1123LabelFmt    = `[a-z0-9]([-a-z0-9]*[a-z0-9])?`
	dns1123LabelRegexp = regexp.MustCompile("^" + dns1123LabelFmt + "$")

	// Common injection patterns to block while allowing hyphens and numbers
	injectionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(\b|\s)(select|insert|update|drop|alter|create|exec|xp_cmdshell|;|/\*|\*/|@@|char\(|--\s)`), // SQL injection
		regexp.MustCompile(`<[\s\/\?]?[^\w\s\/\?\-][^>]*>|<(\w+)[^>]*>|<\/\w+>`),                                            // HTML/XML injection (allows hyphens in tags)
		regexp.MustCompile(`[\$&+,\:;=\?@#|'<>.^\*\[\]()!\/]`),                                                              // Special characters that could be used for injection (allows - and numbers)
	}
)

// ValidateResourceName validates that a name is a valid DNS-1123 label and prevents injection
func ValidateResourceName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > MaxNameLength {
		return fmt.Errorf("name cannot be longer than %d characters", MaxNameLength)
	}
	if !dns1123LabelRegexp.MatchString(name) {
		return fmt.Errorf("name must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character")
	}
	if containsInjectionPatterns(name) {
		return fmt.Errorf("name contains potentially dangerous patterns (allowed: a-z, 0-9, -)")
	}
	return nil
}

// ValidateOrgName validates organization name format and prevents injection
func ValidateOrgName(org string) error {
	if org == "" {
		return fmt.Errorf("organization name cannot be empty")
	}
	if len(org) > MaxNameLength {
		return fmt.Errorf("organization name cannot be longer than %d characters", MaxNameLength)
	}
	if strings.Contains(org, "/") {
		return fmt.Errorf("organization name cannot contain '/'")
	}
	if containsInjectionPatterns(org) {
		return fmt.Errorf("organization name contains potentially dangerous characters or patterns")
	}
	return nil
}

// containsInjectionPatterns checks if the input contains any known injection patterns
func containsInjectionPatterns(input string) bool {
	for _, pattern := range injectionPatterns {
		if pattern.MatchString(input) {
			return true
		}
	}
	return false
}

// sanitizeInput removes any potentially dangerous characters from input
func sanitizeInput(input string) string {
	// First remove any null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove any injection patterns
	for _, pattern := range injectionPatterns {
		input = pattern.ReplaceAllString(input, "")
	}

	return strings.TrimSpace(input)
}

// ValidateBidPrice validates the bid price format and prevents injection
func ValidateBidPrice(price string) error {
	if price == "" {
		return fmt.Errorf("bid price cannot be empty")
	}

	// Check for injection patterns
	if containsInjectionPatterns(price) {
		return fmt.Errorf("bid price contains potentially dangerous characters or patterns")
	}

	// Remove $ prefix if present
	if price[0] == '$' {
		price = price[1:]
	}
	// TODO: Add more specific validation based on your pricing format
	return nil
}

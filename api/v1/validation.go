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
	uuidRegexp         = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

	injectionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(\b|\s)(select|insert|update|drop|alter|create|exec|xp_cmdshell|;|/\*|\*/|@@|char\(|--\s)`), // SQL keywords
		regexp.MustCompile(`<[\s\/\?]?[^\\w\s\/\?\-][^>]*>|<(\w+)[^>]*>|<\/\w+>`),                                           // HTML/XML injection
		regexp.MustCompile(`[\$&+,:;=?@#|'<>^\*\[\]()!\/]`),                                                                 // Special chars excluding dot '.'
	}

	// Regex to match optional $ prefix, digits, optional decimal point with any number of decimals
	currencyRegexp = regexp.MustCompile(`^\$?([0-9]+)(\.[0-9]+)?$`)
)

func isValidUUID(name string) bool {
	return uuidRegexp.MatchString(strings.ToLower(name))
}

// ValidateResourceName validates that a name is a valid DNS-1123 label and prevents injection
func ValidateResourceName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > MaxNameLength {
		return fmt.Errorf("name cannot be longer than %d characters", MaxNameLength)
	}
	if !dns1123LabelRegexp.MatchString(name) && !isValidUUID(name) {
		return fmt.Errorf("name - %s must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character", name)
	}
	if containsInjectionPatterns(name) {
		return fmt.Errorf("name - %s contains potentially dangerous patterns (allowed: a-z, 0-9, -)", name)
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
		return fmt.Errorf("organization name - %s contains potentially dangerous characters or patterns", org)
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

// ValidateBidPrice validates the bid price format and prevents injection
func ValidateBidPrice(price string) error {
	price = strings.TrimSpace(price)
	if price == "" {
		return fmt.Errorf("bid price cannot be empty")
	}

	// Remove $ prefix if present
	priceClean := strings.TrimPrefix(price, "$")

	// Check for dangerous/injection patterns before numeric validation
	if containsInjectionPatterns(priceClean) {
		return fmt.Errorf("bid price - %s contains potentially dangerous characters or patterns", price)
	}

	if !currencyRegexp.MatchString(price) {
		return fmt.Errorf("bid price - %s is not a valid currency format", price)
	}

	// You can optionally parse to float64 here if needed to further validate
	return nil
}

func ValidateServerClass(serverClass ServerClass) error {
	if serverClass.Availability != "available" {
		return fmt.Errorf("server class - %s is not available", serverClass.Name)
	}
	return nil
}

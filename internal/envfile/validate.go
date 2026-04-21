package envfile

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidationError represents a single validation issue found in an EnvMap.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

// Valid returns true if no validation errors were found.
func (r *ValidationResult) Valid() bool {
	return len(r.Errors) == 0
}

// Summary returns a human-readable summary of all validation errors.
func (r *ValidationResult) Summary() string {
	if r.Valid() {
		return "no validation errors"
	}
	var sb strings.Builder
	for _, e := range r.Errors {
		sb.WriteString("  - ")
		sb.WriteString(e.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

// Validate checks an EnvMap for common issues:
//   - Empty keys
//   - Keys containing spaces or invalid characters
//   - Keys that start with a digit
func Validate(env EnvMap) *ValidationResult {
	result := &ValidationResult{}

	for key := range env {
		if key == "" {
			result.Errors = append(result.Errors, ValidationError{
				Key:     key,
				Message: "key must not be empty",
			})
			continue
		}

		if unicode.IsDigit(rune(key[0])) {
			result.Errors = append(result.Errors, ValidationError{
				Key:     key,
				Message: "key must not start with a digit",
			})
		}

		for _, ch := range key {
			if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
				result.Errors = append(result.Errors, ValidationError{
					Key:     key,
					Message: fmt.Sprintf("key contains invalid character %q (only letters, digits, and underscores allowed)", ch),
				})
				break
			}
		}
	}

	return result
}

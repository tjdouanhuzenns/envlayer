package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaField defines an expected environment variable with optional constraints.
type SchemaField struct {
	Key      string
	Required bool
	Pattern  *regexp.Regexp // optional regex the value must match
	AllowedValues []string  // optional allowlist
}

// Schema is a collection of SchemaFields used to validate an EnvMap.
type Schema struct {
	Fields []SchemaField
}

// SchemaError represents a single schema validation failure.
type SchemaError struct {
	Key     string
	Message string
}

func (e SchemaError) Error() string {
	return fmt.Sprintf("schema error for %q: %s", e.Key, e.Message)
}

// ValidateSchema checks an EnvMap against a Schema and returns all violations.
func ValidateSchema(env EnvMap, schema Schema) []SchemaError {
	var errs []SchemaError

	for _, field := range schema.Fields {
		val, exists := env[field.Key]

		if !exists || strings.TrimSpace(val) == "" {
			if field.Required {
				errs = append(errs, SchemaError{Key: field.Key, Message: "required key is missing or empty"})
			}
			continue
		}

		if field.Pattern != nil && !field.Pattern.MatchString(val) {
			errs = append(errs, SchemaError{
				Key:     field.Key,
				Message: fmt.Sprintf("value %q does not match required pattern %s", val, field.Pattern.String()),
			})
		}

		if len(field.AllowedValues) > 0 {
			allowed := false
			for _, av := range field.AllowedValues {
				if av == val {
					allowed = true
					break
				}
			}
			if !allowed {
				errs = append(errs, SchemaError{
					Key:     field.Key,
					Message: fmt.Sprintf("value %q is not in allowed values [%s]", val, strings.Join(field.AllowedValues, ", ")),
				})
			}
		}
	}

	return errs
}

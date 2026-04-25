package envfile

import (
	"regexp"
	"testing"
)

func TestValidateSchema_AllValid(t *testing.T) {
	env := EnvMap{"APP_ENV": "prod", "PORT": "8080"}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "APP_ENV", Required: true, AllowedValues: []string{"dev", "staging", "prod"}},
			{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
		},
	}
	errs := ValidateSchema(env, schema)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	env := EnvMap{"PORT": "8080"}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "APP_ENV", Required: true},
		},
	}
	errs := ValidateSchema(env, schema)
	if len(errs) != 1 || errs[0].Key != "APP_ENV" {
		t.Fatalf("expected one error for APP_ENV, got %v", errs)
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	env := EnvMap{"PORT": "not-a-number"}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
		},
	}
	errs := ValidateSchema(env, schema)
	if len(errs) != 1 {
		t.Fatalf("expected pattern error, got %v", errs)
	}
}

func TestValidateSchema_AllowedValueViolation(t *testing.T) {
	env := EnvMap{"APP_ENV": "local"}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "APP_ENV", Required: true, AllowedValues: []string{"dev", "staging", "prod"}},
		},
	}
	errs := ValidateSchema(env, schema)
	if len(errs) != 1 || errs[0].Key != "APP_ENV" {
		t.Fatalf("expected allowed-value error for APP_ENV, got %v", errs)
	}
}

func TestValidateSchema_OptionalMissingIsOk(t *testing.T) {
	env := EnvMap{}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "LOG_LEVEL", Required: false},
		},
	}
	errs := ValidateSchema(env, schema)
	if len(errs) != 0 {
		t.Fatalf("expected no errors for optional missing key, got %v", errs)
	}
}

func TestValidateSchema_MultipleErrors(t *testing.T) {
	env := EnvMap{"APP_ENV": "unknown"}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "APP_ENV", Required: true, AllowedValues: []string{"dev", "prod"}},
			{Key: "PORT", Required: true},
		},
	}
	errs := ValidateSchema(env, schema)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(errs), errs)
	}
}

func TestSchemaError_String(t *testing.T) {
	e := SchemaError{Key: "FOO", Message: "required key is missing or empty"}
	got := e.Error()
	if got == "" {
		t.Fatal("expected non-empty error string")
	}
}

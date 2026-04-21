package envfile

import (
	"testing"
)

func TestValidate_ValidMap(t *testing.T) {
	env := EnvMap{
		"APP_ENV":       "production",
		"DB_HOST":       "localhost",
		"PORT":          "8080",
		"_PRIVATE_KEY":  "secret",
		"KEY123":        "value",
	}
	result := Validate(env)
	if !result.Valid() {
		t.Errorf("expected valid, got errors:\n%s", result.Summary())
	}
}

func TestValidate_KeyStartsWithDigit(t *testing.T) {
	env := EnvMap{
		"1INVALID": "value",
	}
	result := Validate(env)
	if result.Valid() {
		t.Fatal("expected validation error for key starting with digit")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Key != "1INVALID" {
		t.Errorf("unexpected key in error: %q", result.Errors[0].Key)
	}
}

func TestValidate_KeyWithSpace(t *testing.T) {
	env := EnvMap{
		"BAD KEY": "value",
	}
	result := Validate(env)
	if result.Valid() {
		t.Fatal("expected validation error for key with space")
	}
}

func TestValidate_KeyWithHyphen(t *testing.T) {
	env := EnvMap{
		"BAD-KEY": "value",
	}
	result := Validate(env)
	if result.Valid() {
		t.Fatal("expected validation error for key with hyphen")
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	env := EnvMap{
		"1BAD":    "v",
		"ALSO-BAD": "v",
		"GOOD_KEY": "v",
	}
	result := Validate(env)
	if result.Valid() {
		t.Fatal("expected validation errors")
	}
	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d: %s", len(result.Errors), result.Summary())
	}
}

func TestValidate_EmptyMap(t *testing.T) {
	env := EnvMap{}
	result := Validate(env)
	if !result.Valid() {
		t.Errorf("expected empty map to be valid, got: %s", result.Summary())
	}
}

func TestValidationResult_Summary_NoErrors(t *testing.T) {
	r := &ValidationResult{}
	got := r.Summary()
	if got != "no validation errors" {
		t.Errorf("unexpected summary: %q", got)
	}
}

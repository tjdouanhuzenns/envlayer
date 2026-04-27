package envfile

import (
	"testing"
)

func TestRedact_ExactKey(t *testing.T) {
	env := EnvMap{"DB_PASSWORD": "secret", "APP_HOST": "localhost"}
	out, err := Redact(env, RedactOptions{Keys: []string{"DB_PASSWORD"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", out["DB_PASSWORD"])
	}
	if out["APP_HOST"] != "localhost" {
		t.Errorf("APP_HOST should be unchanged, got %q", out["APP_HOST"])
	}
}

func TestRedact_CaseInsensitiveKey(t *testing.T) {
	env := EnvMap{"DB_PASSWORD": "secret"}
	out, err := Redact(env, RedactOptions{Keys: []string{"db_password"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", out["DB_PASSWORD"])
	}
}

func TestRedact_PatternMatch(t *testing.T) {
	env := EnvMap{"API_SECRET": "abc", "API_KEY": "xyz", "APP_NAME": "envlayer"}
	out, err := Redact(env, RedactOptions{Patterns: []string{"(?i)secret|key"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_SECRET"] != "[REDACTED]" {
		t.Errorf("API_SECRET should be redacted")
	}
	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY should be redacted")
	}
	if out["APP_NAME"] != "envlayer" {
		t.Errorf("APP_NAME should be unchanged")
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	env := EnvMap{"TOKEN": "mytoken"}
	out, err := Redact(env, RedactOptions{Keys: []string{"TOKEN"}, Placeholder: "***"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "***" {
		t.Errorf("expected ***, got %q", out["TOKEN"])
	}
}

func TestRedact_InvalidPattern(t *testing.T) {
	env := EnvMap{"KEY": "val"}
	_, err := Redact(env, RedactOptions{Patterns: []string{"[invalid"}})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	env := EnvMap{"SECRET": "original"}
	_, err := Redact(env, RedactOptions{Keys: []string{"SECRET"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SECRET"] != "original" {
		t.Error("original EnvMap should not be modified")
	}
}

func TestRedact_EmptyOptions(t *testing.T) {
	env := EnvMap{"A": "1", "B": "2"}
	out, err := Redact(env, RedactOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Error("no keys should be redacted with empty options")
	}
}

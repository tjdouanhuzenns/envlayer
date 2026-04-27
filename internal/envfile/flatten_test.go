package envfile

import (
	"testing"
)

func TestFlatten_DotSeparator(t *testing.T) {
	env := EnvMap{"db.host": "localhost", "db.port": "5432"}
	opts := DefaultFlattenOptions()
	got, err := Flatten(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", got["DB_HOST"])
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", got["DB_PORT"])
	}
}

func TestFlatten_DoubleUnderscoreSeparator(t *testing.T) {
	env := EnvMap{"APP__NAME": "envlayer", "APP__ENV": "prod"}
	opts := DefaultFlattenOptions()
	got, err := Flatten(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["APP_NAME"] != "envlayer" {
		t.Errorf("expected APP_NAME=envlayer, got %q", got["APP_NAME"])
	}
	if got["APP_ENV"] != "prod" {
		t.Errorf("expected APP_ENV=prod, got %q", got["APP_ENV"])
	}
}

func TestFlatten_NoSeparatorPassthrough(t *testing.T) {
	env := EnvMap{"PLAIN_KEY": "value"}
	opts := DefaultFlattenOptions()
	got, err := Flatten(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["PLAIN_KEY"] != "value" {
		t.Errorf("expected PLAIN_KEY=value, got %q", got["PLAIN_KEY"])
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	env := EnvMap{"db.host": "localhost"}
	opts := FlattenOptions{Separator: "_", UpperCase: true, Prefix: "APP_"}
	got, err := Flatten(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["APP_DB_HOST"] != "localhost" {
		t.Errorf("expected APP_DB_HOST=localhost, got %q", got["APP_DB_HOST"])
	}
}

func TestFlatten_CollisionReturnsError(t *testing.T) {
	// Both keys normalise to DB_HOST after flattening.
	env := EnvMap{"db.host": "a", "db__host": "b"}
	opts := DefaultFlattenOptions()
	_, err := Flatten(env, opts)
	if err == nil {
		t.Fatal("expected collision error, got nil")
	}
}

func TestFlatten_NilInputReturnsError(t *testing.T) {
	_, err := Flatten(nil, DefaultFlattenOptions())
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestFlatten_LowerCasePreserved(t *testing.T) {
	env := EnvMap{"db.host": "localhost"}
	opts := FlattenOptions{Separator: "_", UpperCase: false}
	got, err := Flatten(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["db_host"] != "localhost" {
		t.Errorf("expected db_host=localhost, got %q", got["db_host"])
	}
}

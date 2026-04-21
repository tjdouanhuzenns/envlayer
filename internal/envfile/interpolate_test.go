package envfile

import (
	"os"
	"testing"
)

func TestInterpolate_BasicReference(t *testing.T) {
	env := EnvMap{"HOST": "localhost", "URL": "http://${HOST}:8080"}
	result, err := Interpolate(env, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result["URL"]; got != "http://localhost:8080" {
		t.Errorf("URL = %q, want %q", got, "http://localhost:8080")
	}
}

func TestInterpolate_DollarNoBraces(t *testing.T) {
	env := EnvMap{"NAME": "world", "GREETING": "hello $NAME"}
	result, err := Interpolate(env, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result["GREETING"]; got != "hello world" {
		t.Errorf("GREETING = %q, want %q", got, "hello world")
	}
}

func TestInterpolate_MissingKeyEmptyString(t *testing.T) {
	env := EnvMap{"VAL": "prefix_${MISSING}_suffix"}
	result, err := Interpolate(env, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result["VAL"]; got != "prefix__suffix" {
		t.Errorf("VAL = %q, want %q", got, "prefix__suffix")
	}
}

func TestInterpolate_FailOnMissing(t *testing.T) {
	env := EnvMap{"VAL": "${UNDEFINED}"}
	_, err := Interpolate(env, InterpolateOptions{FailOnMissing: true})
	if err == nil {
		t.Fatal("expected error for undefined variable, got nil")
	}
}

func TestInterpolate_FallbackToOS(t *testing.T) {
	os.Setenv("OS_VAR", "from-os")
	defer os.Unsetenv("OS_VAR")

	env := EnvMap{"COMBINED": "val=${OS_VAR}"}
	result, err := Interpolate(env, InterpolateOptions{FallbackToOS: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result["COMBINED"]; got != "val=from-os" {
		t.Errorf("COMBINED = %q, want %q", got, "val=from-os")
	}
}

func TestInterpolate_NoReferences(t *testing.T) {
	env := EnvMap{"PLAIN": "no-refs-here", "NUM": "42"}
	result, err := Interpolate(env, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["PLAIN"] != "no-refs-here" || result["NUM"] != "42" {
		t.Errorf("values should be unchanged: %v", result)
	}
}

func TestInterpolate_EmptyMap(t *testing.T) {
	result, err := Interpolate(EnvMap{}, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

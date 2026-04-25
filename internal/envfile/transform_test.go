package envfile

import (
	"errors"
	"strings"
	"testing"
)

func TestTransform_UpperCase(t *testing.T) {
	env := EnvMap{"FOO": "hello", "BAR": "world"}
	out, results, err := Transform(env, BuiltinTransforms.UpperCase, TransformOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "HELLO" || out["BAR"] != "WORLD" {
		t.Errorf("expected upper-cased values, got %v", out)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestTransform_LimitedKeys(t *testing.T) {
	env := EnvMap{"FOO": "hello", "BAR": "world"}
	opts := TransformOptions{Keys: []string{"FOO"}}
	out, _, err := Transform(env, BuiltinTransforms.UpperCase, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "HELLO" {
		t.Errorf("FOO should be upper-cased")
	}
	if out["BAR"] != "world" {
		t.Errorf("BAR should be unchanged, got %q", out["BAR"])
	}
}

func TestTransform_TrimSpace(t *testing.T) {
	env := EnvMap{"KEY": "  spaced  "}
	out, _, err := Transform(env, BuiltinTransforms.TrimSpace, TransformOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "spaced" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestTransform_PrefixKeys(t *testing.T) {
	env := EnvMap{"NAME": "alice"}
	out, _, _ := Transform(env, BuiltinTransforms.PrefixKeys("APP_"), TransformOptions{})
	if out["NAME"] != "APP_alice" {
		t.Errorf("expected prefixed value, got %q", out["NAME"])
	}
}

func TestTransform_ErrorPreservesOriginal(t *testing.T) {
	env := EnvMap{"KEY": "original"}
	failFn := func(k, v string) (string, error) {
		return "", errors.New("oops")
	}
	out, results, err := Transform(env, failFn, TransformOptions{})
	if err != nil {
		t.Fatalf("unexpected top-level error when FailOnError=false")
	}
	if out["KEY"] != "original" {
		t.Errorf("original value should be preserved on error")
	}
	if results[0].Err == nil {
		t.Errorf("expected error in result")
	}
}

func TestTransform_FailOnError(t *testing.T) {
	env := EnvMap{"A": "val"}
	failFn := func(k, v string) (string, error) {
		return "", errors.New("forced failure")
	}
	_, _, err := Transform(env, failFn, TransformOptions{FailOnError: true})
	if err == nil {
		t.Error("expected error when FailOnError=true")
	}
	if !strings.Contains(err.Error(), "A") {
		t.Errorf("error should mention key name, got: %v", err)
	}
}

func TestTransform_NilEnv(t *testing.T) {
	out, results, err := Transform(nil, BuiltinTransforms.UpperCase, TransformOptions{})
	if err != nil {
		t.Fatalf("unexpected error on nil env")
	}
	if len(out) != 0 || len(results) != 0 {
		t.Errorf("expected empty output for nil env")
	}
}

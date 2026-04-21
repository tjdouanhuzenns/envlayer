package envfile

import (
	"strings"
	"testing"
)

func TestMask_DefaultPatterns(t *testing.T) {
	env := EnvMap{
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abc123",
		"APP_NAME":    "myapp",
		"SECRET_KEY":  "topsecret",
	}
	masked := Mask(env, nil)

	for _, k := range []string{"DB_PASSWORD", "API_KEY", "SECRET_KEY"} {
		if masked[k] == env[k] {
			t.Errorf("expected %s to be masked, got %q", k, masked[k])
		}
		if !strings.Contains(masked[k], "*") {
			t.Errorf("expected %s masked value to contain '*', got %q", k, masked[k])
		}
	}
	if masked["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", masked["APP_NAME"])
	}
}

func TestMask_ExactKey(t *testing.T) {
	env := EnvMap{"MY_CUSTOM": "sensitive", "OTHER": "safe"}
	masked := Mask(env, &MaskOptions{Keys: []string{"MY_CUSTOM"}})
	if masked["MY_CUSTOM"] == "sensitive" {
		t.Error("expected MY_CUSTOM to be masked")
	}
	if masked["OTHER"] != "safe" {
		t.Errorf("expected OTHER unchanged, got %q", masked["OTHER"])
	}
}

func TestMask_VisibleChars(t *testing.T) {
	env := EnvMap{"API_TOKEN": "abcdefghij"}
	masked := Mask(env, &MaskOptions{VisibleChars: 4})
	val := masked["API_TOKEN"]
	if !strings.HasSuffix(val, "ghij") {
		t.Errorf("expected masked value to end with 'ghij', got %q", val)
	}
	if !strings.HasPrefix(val, "******") {
		t.Errorf("expected masked value to start with '******', got %q", val)
	}
}

func TestMask_EmptyValue(t *testing.T) {
	env := EnvMap{"DB_PASSWORD": ""}
	masked := Mask(env, nil)
	if masked["DB_PASSWORD"] != "" {
		t.Errorf("expected empty value to stay empty, got %q", masked["DB_PASSWORD"])
	}
}

func TestMask_CustomMaskChar(t *testing.T) {
	env := EnvMap{"SECRET": "mysecret"}
	masked := Mask(env, &MaskOptions{MaskChar: "#"})
	if !strings.Contains(masked["SECRET"], "#") {
		t.Errorf("expected custom mask char '#', got %q", masked["SECRET"])
	}
}

func TestMask_DoesNotMutateOriginal(t *testing.T) {
	env := EnvMap{"DB_PASSWORD": "original"}
	_ = Mask(env, nil)
	if env["DB_PASSWORD"] != "original" {
		t.Error("Mask mutated the original EnvMap")
	}
}

func TestMask_NilOpts_UsesDefaults(t *testing.T) {
	env := EnvMap{"TOKEN": "tok_live_xyz", "PORT": "8080"}
	masked := Mask(env, nil)
	if masked["TOKEN"] == "tok_live_xyz" {
		t.Error("expected TOKEN to be masked with nil opts")
	}
	if masked["PORT"] != "8080" {
		t.Errorf("expected PORT unchanged, got %q", masked["PORT"])
	}
}

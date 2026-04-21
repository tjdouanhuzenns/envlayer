package envfile_test

import (
	"strings"
	"testing"

	"github.com/user/envlayer/internal/envfile"
)

// TestMask_WithExport verifies that Mask integrates cleanly with Export,
// so masked values are written to output formats without leaking secrets.
func TestMask_WithExport(t *testing.T) {
	env := envfile.EnvMap{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "key-abc-123",
	}

	masked := envfile.Mask(env, nil)

	out, err := envfile.Export(masked, envfile.FormatDotenv)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if strings.Contains(out, "supersecret") {
		t.Error("exported output should not contain original DB_PASSWORD value")
	}
	if strings.Contains(out, "key-abc-123") {
		t.Error("exported output should not contain original API_KEY value")
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Error("exported output should still contain APP_NAME")
	}
}

// TestMask_WithMerge verifies masking after a merge preserves non-sensitive keys.
func TestMask_WithMerge(t *testing.T) {
	base := envfile.EnvMap{
		"APP_ENV":  "production",
		"DB_PASS":  "base-pass",
	}
	overlay := envfile.EnvMap{
		"DB_PASS":  "overlay-pass",
		"NEW_TOKEN": "tok_xyz",
	}

	merged := envfile.Merge(base, overlay)
	masked := envfile.Mask(merged, nil)

	if masked["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged, got %q", masked["APP_ENV"])
	}
	if masked["DB_PASS"] == "overlay-pass" {
		t.Error("DB_PASS should be masked after merge")
	}
	if masked["NEW_TOKEN"] == "tok_xyz" {
		t.Error("NEW_TOKEN should be masked after merge")
	}
}

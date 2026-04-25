package envfile

import (
	"strings"
	"testing"
)

// TestTransform_ThenExport verifies that a transformed map can be exported.
func TestTransform_ThenExport(t *testing.T) {
	base := EnvMap{
		"DB_HOST": "  localhost  ",
		"DB_PORT": "5432",
	}

	trimmed, _, err := Transform(base, BuiltinTransforms.TrimSpace, TransformOptions{})
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	out, err := Export(trimmed, FormatDotenv)
	if err != nil {
		t.Fatalf("export error: %v", err)
	}

	if strings.Contains(out, "  localhost  ") {
		t.Errorf("expected trimmed value in export, got raw spaces")
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected 'localhost' in export output")
	}
}

// TestTransform_AfterMerge verifies transform can be applied post-merge.
func TestTransform_AfterMerge(t *testing.T) {
	base := EnvMap{"APP_ENV": "development", "LOG_LEVEL": "debug"}
	layer := EnvMap{"APP_ENV": "production"}

	merged := Merge(base, layer)

	upper, _, err := Transform(merged, BuiltinTransforms.UpperCase, TransformOptions{
		Keys: []string{"APP_ENV"},
	})
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	if upper["APP_ENV"] != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %q", upper["APP_ENV"])
	}
	if upper["LOG_LEVEL"] != "debug" {
		t.Errorf("LOG_LEVEL should be unchanged, got %q", upper["LOG_LEVEL"])
	}
}

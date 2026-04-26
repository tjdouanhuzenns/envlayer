package envfile_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envlayer/internal/envfile"
)

func TestPipeline_SingleStep(t *testing.T) {
	base := envfile.NewEnvMap(map[string]string{
		"APP_NAME": "myapp",
		"DEBUG":    "true",
	})

	p := envfile.NewPipeline(base)
	p.Step("uppercase keys", func(m map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(m))
		for k, v := range m {
			out[strings.ToUpper(k)] = v
		}
		return out, nil
	})

	result, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %q", result["APP_NAME"])
	}
}

func TestPipeline_MultipleSteps(t *testing.T) {
	base := envfile.NewEnvMap(map[string]string{
		"host": "  localhost  ",
		"port": "8080",
	})

	p := envfile.NewPipeline(base)

	// Step 1: trim whitespace from values
	p.Step("trim values", func(m map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(m))
		for k, v := range m {
			out[k] = strings.TrimSpace(v)
		}
		return out, nil
	})

	// Step 2: uppercase keys
	p.Step("uppercase keys", func(m map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(m))
		for k, v := range m {
			out[strings.ToUpper(k)] = v
		}
		return out, nil
	})

	result, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", result["HOST"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", result["PORT"])
	}
}

func TestPipeline_StepError(t *testing.T) {
	base := envfile.NewEnvMap(map[string]string{
		"KEY": "value",
	})

	p := envfile.NewPipeline(base)
	p.Step("failing step", func(m map[string]string) (map[string]string, error) {
		return nil, fmt.Errorf("step failed")
	})
	p.Step("should not run", func(m map[string]string) (map[string]string, error) {
		t.Error("this step should not have been called")
		return m, nil
	})

	_, err := p.Run()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "step failed") {
		t.Errorf("expected error to contain 'step failed', got: %v", err)
	}
}

func TestPipeline_EmptySteps(t *testing.T) {
	base := envfile.NewEnvMap(map[string]string{
		"FOO": "bar",
	})

	p := envfile.NewPipeline(base)
	result, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", result["FOO"])
	}
}

func TestPipeline_IsolatesInput(t *testing.T) {
	original := map[string]string{
		"KEY": "original",
	}
	base := envfile.NewEnvMap(original)

	p := envfile.NewPipeline(base)
	p.Step("mutate", func(m map[string]string) (map[string]string, error) {
		m["KEY"] = "mutated"
		return m, nil
	})

	_, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// original map should not be mutated
	if original["KEY"] != "original" {
		t.Errorf("pipeline mutated the original input map")
	}
}

func TestPipeline_WithRealTransform(t *testing.T) {
	base := envfile.NewEnvMap(map[string]string{
		"db_host": "localhost",
		"db_port": "5432",
	})

	opts := envfile.TransformOptions{
		KeyPrefix: "APP_",
	}

	p := envfile.NewPipeline(base)
	p.Step("prefix keys", func(m map[string]string) (map[string]string, error) {
		return envfile.Transform(m, opts)
	})

	result, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["APP_db_host"]; !ok {
		t.Errorf("expected key APP_db_host in result, got keys: %v", result)
	}
}

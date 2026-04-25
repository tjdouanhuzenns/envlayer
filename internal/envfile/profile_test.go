package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeProfileEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeProfileEnv: %v", err)
	}
	return p
}

func TestProfileRegistry_RegisterAndGet(t *testing.T) {
	reg := NewProfileRegistry()
	p := &Profile{Name: "dev", Layers: []string{".env"}}
	if err := reg.Register(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("dev")
	if !ok || got.Name != "dev" {
		t.Errorf("expected profile dev, got %v", got)
	}
}

func TestProfileRegistry_RegisterNilError(t *testing.T) {
	reg := NewProfileRegistry()
	if err := reg.Register(nil); err == nil {
		t.Error("expected error for nil profile")
	}
}

func TestProfileRegistry_RegisterEmptyNameError(t *testing.T) {
	reg := NewProfileRegistry()
	if err := reg.Register(&Profile{Name: ""}); err == nil {
		t.Error("expected error for empty profile name")
	}
}

func TestProfileRegistry_List(t *testing.T) {
	reg := NewProfileRegistry()
	for _, n := range []string{"prod", "dev", "staging"} {
		_ = reg.Register(&Profile{Name: n})
	}
	list := reg.List()
	if list[0] != "dev" || list[1] != "prod" || list[2] != "staging" {
		t.Errorf("unexpected order: %v", list)
	}
}

func TestProfileRegistry_ResolveNotFound(t *testing.T) {
	reg := NewProfileRegistry()
	_, err := reg.Resolve("missing")
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestProfileRegistry_ResolveEmptyLayers(t *testing.T) {
	reg := NewProfileRegistry()
	_ = reg.Register(&Profile{Name: "empty"})
	env, err := reg.Resolve("empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 0 {
		t.Errorf("expected empty map, got %v", env)
	}
}

func TestProfileRegistry_ResolveMergesLayers(t *testing.T) {
	dir := t.TempDir()
	base := writeProfileEnv(t, dir, ".env", "APP=base\nDEBUG=false\n")
	over := writeProfileEnv(t, dir, ".env.staging", "DEBUG=true\nREGION=us-east\n")

	reg := NewProfileRegistry()
	_ = reg.Register(&Profile{Name: "staging", Layers: []string{base, over}})

	env, err := reg.Resolve("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP"] != "base" {
		t.Errorf("APP: want base, got %s", env["APP"])
	}
	if env["DEBUG"] != "true" {
		t.Errorf("DEBUG: want true, got %s", env["DEBUG"])
	}
	if env["REGION"] != "us-east" {
		t.Errorf("REGION: want us-east, got %s", env["REGION"])
	}
}

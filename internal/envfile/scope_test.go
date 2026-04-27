package envfile

import (
	"testing"
)

func TestScopeRegistry_RegisterAndGet(t *testing.T) {
	reg := NewScopeRegistry()
	env := EnvMap{"KEY": "value"}
	if err := reg.Register("dev", env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("dev")
	if !ok {
		t.Fatal("expected to find scope 'dev'")
	}
	if got["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", got["KEY"])
	}
}

func TestScopeRegistry_RegisterEmptyNameError(t *testing.T) {
	reg := NewScopeRegistry()
	if err := reg.Register("", EnvMap{}); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestScopeRegistry_RegisterDuplicateError(t *testing.T) {
	reg := NewScopeRegistry()
	_ = reg.Register("prod", EnvMap{"A": "1"})
	if err := reg.Register("prod", EnvMap{"B": "2"}); err == nil {
		t.Fatal("expected error for duplicate scope name")
	}
}

func TestScopeRegistry_GetMissing(t *testing.T) {
	reg := NewScopeRegistry()
	_, ok := reg.Get("staging")
	if ok {
		t.Fatal("expected missing scope to return false")
	}
}

func TestScopeRegistry_List(t *testing.T) {
	reg := NewScopeRegistry()
	_ = reg.Register("base", EnvMap{})
	_ = reg.Register("dev", EnvMap{})
	_ = reg.Register("prod", EnvMap{})
	names := reg.List()
	if len(names) != 3 {
		t.Fatalf("expected 3 scopes, got %d", len(names))
	}
	if names[0] != "base" || names[1] != "dev" || names[2] != "prod" {
		t.Errorf("unexpected order: %v", names)
	}
}

func TestScopeRegistry_Resolve_Override(t *testing.T) {
	reg := NewScopeRegistry()
	_ = reg.Register("base", EnvMap{"HOST": "localhost", "PORT": "5432"})
	_ = reg.Register("prod", EnvMap{"HOST": "db.prod.example.com"})

	result, err := reg.Resolve("base", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["HOST"] != "db.prod.example.com" {
		t.Errorf("expected prod HOST, got %q", result["HOST"])
	}
	if result["PORT"] != "5432" {
		t.Errorf("expected base PORT, got %q", result["PORT"])
	}
}

func TestScopeRegistry_Resolve_MissingScope(t *testing.T) {
	reg := NewScopeRegistry()
	_ = reg.Register("base", EnvMap{"A": "1"})
	_, err := reg.Resolve("base", "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing scope")
	}
}

func TestScopeRegistry_Resolve_IsolatesInput(t *testing.T) {
	reg := NewScopeRegistry()
	base := EnvMap{"KEY": "original"}
	_ = reg.Register("base", base)
	result, _ := reg.Resolve("base")
	result["KEY"] = "mutated"
	got, _ := reg.Get("base")
	if got["KEY"] != "original" {
		t.Error("Resolve should not mutate the registered scope")
	}
}

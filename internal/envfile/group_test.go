package envfile

import (
	"testing"
)

func TestGroupRegistry_RegisterAndGet(t *testing.T) {
	r := NewGroupRegistry()
	if err := r.Register("db", []string{"DB_HOST", "DB_PORT"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g, err := r.Get("db")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if g.Name != "db" || len(g.Keys) != 2 {
		t.Errorf("unexpected group: %+v", g)
	}
}

func TestGroupRegistry_RegisterEmptyNameError(t *testing.T) {
	r := NewGroupRegistry()
	if err := r.Register("", []string{"KEY"}); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestGroupRegistry_RegisterNoKeysError(t *testing.T) {
	r := NewGroupRegistry()
	if err := r.Register("empty", []string{}); err == nil {
		t.Error("expected error for empty keys slice")
	}
}

func TestGroupRegistry_RegisterDuplicateError(t *testing.T) {
	r := NewGroupRegistry()
	_ = r.Register("db", []string{"DB_HOST"})
	if err := r.Register("db", []string{"DB_PORT"}); err == nil {
		t.Error("expected error for duplicate group name")
	}
}

func TestGroupRegistry_GetMissing(t *testing.T) {
	r := NewGroupRegistry()
	_, err := r.Get("nonexistent")
	if err == nil {
		t.Error("expected error for missing group")
	}
}

func TestGroupRegistry_List(t *testing.T) {
	r := NewGroupRegistry()
	_ = r.Register("web", []string{"PORT"})
	_ = r.Register("db", []string{"DB_HOST"})
	list := r.List()
	if len(list) != 2 || list[0] != "db" || list[1] != "web" {
		t.Errorf("unexpected list order: %v", list)
	}
}

func TestGroupRegistry_Extract_Present(t *testing.T) {
	r := NewGroupRegistry()
	_ = r.Register("db", []string{"DB_HOST", "DB_PORT"})
	src := EnvMap{"DB_HOST": "localhost", "DB_PORT": "5432", "APP_KEY": "secret"}
	out, err := r.Extract("db", src, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["APP_KEY"]; ok {
		t.Error("APP_KEY should not be in extracted group")
	}
}

func TestGroupRegistry_Extract_SkipMissing(t *testing.T) {
	r := NewGroupRegistry()
	_ = r.Register("db", []string{"DB_HOST", "DB_PASS"})
	src := EnvMap{"DB_HOST": "localhost"}
	out, err := r.Extract("db", src, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestGroupRegistry_Extract_StrictMissing(t *testing.T) {
	r := NewGroupRegistry()
	_ = r.Register("db", []string{"DB_HOST", "DB_PASS"})
	src := EnvMap{"DB_HOST": "localhost"}
	_, err := r.Extract("db", src, true)
	if err == nil {
		t.Error("expected error for missing key in strict mode")
	}
}

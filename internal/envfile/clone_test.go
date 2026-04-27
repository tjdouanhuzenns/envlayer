package envfile

import (
	"testing"
)

func TestClone_NilSource(t *testing.T) {
	out, err := Clone(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestClone_AllKeys(t *testing.T) {
	src := EnvMap{"A": "1", "B": "2", "C": "3"}
	out, err := Clone(src, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 keys, got %d", len(out))
	}
	// Verify isolation — mutating src should not affect clone.
	src["A"] = "changed"
	if out["A"] != "1" {
		t.Errorf("clone was not isolated from source mutation")
	}
}

func TestClone_WithPrefix(t *testing.T) {
	src := EnvMap{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_HOST": "db"}
	out, err := Clone(src, &CloneOptions{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("DB_HOST should not be present")
	}
}

func TestClone_StripPrefix(t *testing.T) {
	src := EnvMap{"APP_HOST": "localhost", "APP_PORT": "8080", "OTHER": "x"}
	out, err := Clone(src, &CloneOptions{Prefix: "APP_", StripPrefix: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", out["HOST"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", out["PORT"])
	}
	if _, ok := out["OTHER"]; ok {
		t.Error("OTHER should not be present after prefix filter")
	}
}

func TestClone_ExcludeKeys(t *testing.T) {
	src := EnvMap{"A": "1", "SECRET": "s3cr3t", "B": "2"}
	out, err := Clone(src, &CloneOptions{ExcludeKeys: []string{"SECRET"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["SECRET"]; ok {
		t.Error("SECRET should be excluded")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestClone_PrefixEqualsKey(t *testing.T) {
	// A key that exactly equals the prefix should be skipped when StripPrefix is true.
	src := EnvMap{"APP_": "bare", "APP_NAME": "envlayer"}
	out, err := Clone(src, &CloneOptions{Prefix: "APP_", StripPrefix: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out[""]; ok {
		t.Error("empty key should not appear in output")
	}
	if out["NAME"] != "envlayer" {
		t.Errorf("expected NAME=envlayer, got %q", out["NAME"])
	}
}

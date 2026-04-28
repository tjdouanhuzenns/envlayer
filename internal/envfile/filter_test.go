package envfile

import (
	"testing"
)

func TestFilter_ByPrefix(t *testing.T) {
	env := EnvMap{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_HOST": "db"}
	got, err := Filter(env, FilterOptions{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
	if _, ok := got["DB_HOST"]; ok {
		t.Error("DB_HOST should have been filtered out")
	}
}

func TestFilter_BySuffix(t *testing.T) {
	env := EnvMap{"APP_HOST": "localhost", "DB_HOST": "db", "APP_PORT": "8080"}
	got, err := Filter(env, FilterOptions{Suffix: "_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
	if _, ok := got["APP_PORT"]; ok {
		t.Error("APP_PORT should have been filtered out")
	}
}

func TestFilter_ByPattern(t *testing.T) {
	env := EnvMap{"SECRET_KEY": "s", "API_KEY": "a", "HOST": "h"}
	got, err := Filter(env, FilterOptions{Pattern: "_KEY$"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	env := EnvMap{"KEY": "val"}
	_, err := Filter(env, FilterOptions{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestFilter_ByKeys(t *testing.T) {
	env := EnvMap{"A": "1", "B": "2", "C": "3"}
	got, err := Filter(env, FilterOptions{Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
	if _, ok := got["B"]; ok {
		t.Error("B should have been filtered out")
	}
}

func TestFilter_Invert(t *testing.T) {
	env := EnvMap{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_HOST": "db"}
	got, err := Filter(env, FilterOptions{Prefix: "APP_", Invert: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 key, got %d", len(got))
	}
	if _, ok := got["DB_HOST"]; !ok {
		t.Error("DB_HOST should be present after invert")
	}
}

func TestFilter_EmptyOptions(t *testing.T) {
	env := EnvMap{"A": "1", "B": "2"}
	got, err := Filter(env, FilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(env) {
		t.Fatalf("expected all %d keys, got %d", len(env), len(got))
	}
}

func TestFilter_DoesNotMutateOriginal(t *testing.T) {
	env := EnvMap{"APP_X": "1", "DB_Y": "2"}
	_, err := Filter(env, FilterOptions{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 2 {
		t.Error("original map was mutated")
	}
}

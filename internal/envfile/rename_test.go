package envfile

import (
	"testing"
)

func TestRename_BasicMapping(t *testing.T) {
	env := EnvMap{"OLD_KEY": "value", "KEEP": "same"}
	opts := DefaultRenameOptions()
	opts.Mapping = map[string]string{"OLD_KEY": "NEW_KEY"}

	result, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", result["NEW_KEY"])
	}
	if _, ok := result["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be dropped")
	}
	if result["KEEP"] != "same" {
		t.Errorf("expected KEEP=same, got %q", result["KEEP"])
	}
}

func TestRename_KeepOriginal(t *testing.T) {
	env := EnvMap{"SRC": "hello"}
	opts := RenameOptions{
		Mapping:      map[string]string{"SRC": "DST"},
		DropOriginal: false,
	}

	result, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["SRC"] != "hello" {
		t.Errorf("expected SRC to be preserved")
	}
	if result["DST"] != "hello" {
		t.Errorf("expected DST=hello")
	}
}

func TestRename_MissingKeySkipped(t *testing.T) {
	env := EnvMap{"EXISTING": "val"}
	opts := DefaultRenameOptions()
	opts.Mapping = map[string]string{"MISSING": "NEW"}

	result, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["NEW"]; ok {
		t.Error("expected NEW to not be present")
	}
}

func TestRename_FailOnMissing(t *testing.T) {
	env := EnvMap{"A": "1"}
	opts := RenameOptions{
		Mapping:       map[string]string{"GHOST": "SPIRIT"},
		FailOnMissing: true,
		DropOriginal:  true,
	}

	_, err := Rename(env, opts)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRename_NilMapError(t *testing.T) {
	_, err := Rename(nil, DefaultRenameOptions())
	if err == nil {
		t.Fatal("expected error for nil map")
	}
}

func TestRename_EmptyMappingError(t *testing.T) {
	env := EnvMap{"A": "1"}
	opts := RenameOptions{Mapping: map[string]string{}}
	_, err := Rename(env, opts)
	if err == nil {
		t.Fatal("expected error for empty mapping")
	}
}

func TestRename_EmptyTargetKeyError(t *testing.T) {
	env := EnvMap{"A": "1"}
	opts := RenameOptions{
		Mapping:      map[string]string{"A": ""},
		DropOriginal: true,
	}
	_, err := Rename(env, opts)
	if err == nil {
		t.Fatal("expected error for empty target key")
	}
}

package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeEnvFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
}

func TestResolveFiles_AllPresent(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env", "BASE=1\n")
	writeEnvFile(t, dir, ".env.dev", "ENV=dev\n")
	writeEnvFile(t, dir, ".env.dev.local", "LOCAL=true\n")

	files, err := ResolveFiles(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(files))
	}
}

func TestResolveFiles_OnlyBase(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env", "BASE=1\n")

	files, err := ResolveFiles(dir, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
}

func TestResolveFiles_NoneFound(t *testing.T) {
	dir := t.TempDir()
	_, err := ResolveFiles(dir, "staging")
	if err == nil {
		t.Fatal("expected error when no files found")
	}
}

func TestResolveAndMerge_Override(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env", "HOST=localhost\nPORT=5432\n")
	writeEnvFile(t, dir, ".env.dev", "HOST=devhost\n")
	writeEnvFile(t, dir, ".env.dev.local", "PORT=9999\n")

	result, err := ResolveAndMerge(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["HOST"] != "devhost" {
		t.Errorf("expected HOST=devhost, got %s", result["HOST"])
	}
	if result["PORT"] != "9999" {
		t.Errorf("expected PORT=9999, got %s", result["PORT"])
	}
}

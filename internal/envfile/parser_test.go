package envfile

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envlayer-*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nPORT=8080\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Values["APP_ENV"] != "production" {
		t.Errorf("expected production, got %s", env.Values["APP_ENV"])
	}
	if env.Values["PORT"] != "8080" {
		t.Errorf("expected 8080, got %s", env.Values["PORT"])
	}
}

func TestParseFile_Comments(t *testing.T) {
	path := writeTempEnv(t, "# comment\nKEY=value\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Keys) != 1 || env.Values["KEY"] != "value" {
		t.Errorf("unexpected result: %+v", env)
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret"` + "\n" + `TOKEN='abc123'` + "\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Values["SECRET"] != "my secret" {
		t.Errorf("expected 'my secret', got %q", env.Values["SECRET"])
	}
	if env.Values["TOKEN"] != "abc123" {
		t.Errorf("expected 'abc123', got %q", env.Values["TOKEN"])
	}
}

func TestParseFile_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")
	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for invalid line")
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

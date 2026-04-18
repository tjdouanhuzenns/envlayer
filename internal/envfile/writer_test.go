package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteString_Basic(t *testing.T) {
	env := EnvMap{"FOO": "bar", "BAZ": "qux"}
	out := WriteString(env)
	// keys sorted: BAZ, FOO
	expected := "BAZ=qux\nFOO=bar\n"
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestWriteString_QuotesSpaces(t *testing.T) {
	env := EnvMap{"MSG": "hello world"}
	out := WriteString(env)
	expected := `MSG="hello world"` + "\n"
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestWriteString_EmptyValue(t *testing.T) {
	env := EnvMap{"EMPTY": ""}
	out := WriteString(env)
	if out != "EMPTY=\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestWriteFile_RoundTrip(t *testing.T) {
	original := EnvMap{
		"HOST": "localhost",
		"PORT": "5432",
		"DSN":  "postgres://user:pass@host/db",
	}

	tmp := filepath.Join(t.TempDir(), "out.env")
	if err := WriteFile(tmp, original); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	parsed, err := ParseFile(tmp)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}

	for k, v := range original {
		if parsed[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, parsed[k])
		}
	}
}

func TestWriteFile_CreatesFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "new.env")
	env := EnvMap{"KEY": "value"}
	if err := WriteFile(tmp, env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(tmp); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

package envfile

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExport_Dotenv(t *testing.T) {
	env := EnvMap{"HOST": "localhost", "PORT": "5432"}
	var buf bytes.Buffer
	if err := Export(&buf, env, FormatDotenv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "HOST=localhost") {
		t.Errorf("expected HOST=localhost in output: %s", out)
	}
	if !strings.Contains(out, "PORT=5432") {
		t.Errorf("expected PORT=5432 in output: %s", out)
	}
}

func TestExport_ExportFormat(t *testing.T) {
	env := EnvMap{"APP": "myapp"}
	var buf bytes.Buffer
	if err := Export(&buf, env, FormatExport); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(buf.String(), "export APP=myapp") {
		t.Errorf("expected 'export APP=myapp', got: %s", buf.String())
	}
}

func TestExport_JSON(t *testing.T) {
	env := EnvMap{"KEY": "val"}
	var buf bytes.Buffer
	if err := Export(&buf, env, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"KEY": "val"`) {
		t.Errorf("expected JSON key-value pair, got: %s", out)
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	env := EnvMap{"X": "1"}
	var buf bytes.Buffer
	err := Export(&buf, env, ExportFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestExport_QuotedSpaces(t *testing.T) {
	env := EnvMap{"MSG": "hello world"}
	var buf bytes.Buffer
	_ = Export(&buf, env, FormatDotenv)
	if !strings.Contains(buf.String(), `"hello world"`) {
		t.Errorf("expected quoted value, got: %s", buf.String())
	}
}

func TestExportToFile_Creates(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")
	env := EnvMap{"HELLO": "world"}
	if err := ExportToFile(path, env, FormatDotenv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if !strings.Contains(string(data), "HELLO=world") {
		t.Errorf("expected HELLO=world in file, got: %s", string(data))
	}
}

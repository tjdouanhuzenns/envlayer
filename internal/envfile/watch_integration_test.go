package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestWatch_ThenDiff verifies that a WatchEvent's Env can be diffed against
// a previously captured snapshot to surface meaningful changes.
func TestWatch_ThenDiff(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")

	initial := "HOST=localhost\nPORT=5432\nDEBUG=false\n"
	if err := os.WriteFile(p, []byte(initial), 0644); err != nil {
		t.Fatal(err)
	}

	baseEnv, err := ParseFile(p)
	if err != nil {
		t.Fatal(err)
	}

	w := NewWatcher(20*time.Millisecond, p)
	w.Start()
	defer w.Stop()

	time.Sleep(40 * time.Millisecond)

	updated := "HOST=localhost\nPORT=5433\nDEBUG=true\nNEW_KEY=hello\n"
	if err := os.WriteFile(p, []byte(updated), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case ev := <-w.Events:
		diffs := Diff(baseEnv, ev.Env)
		statuses := map[string]string{}
		for _, d := range diffs {
			statuses[d.Key] = d.Status
		}
		if statuses["PORT"] != "changed" {
			t.Errorf("expected PORT changed, got %q", statuses["PORT"])
		}
		if statuses["DEBUG"] != "changed" {
			t.Errorf("expected DEBUG changed, got %q", statuses["DEBUG"])
		}
		if statuses["NEW_KEY"] != "added" {
			t.Errorf("expected NEW_KEY added, got %q", statuses["NEW_KEY"])
		}
		if statuses["HOST"] != "unchanged" {
			t.Errorf("expected HOST unchanged, got %q", statuses["HOST"])
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for watch event")
	}
}

// TestWatch_ThenMerge verifies that an updated env from a WatchEvent can be
// merged over a base map to produce the final resolved state.
func TestWatch_ThenMerge(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "override.env")

	if err := os.WriteFile(p, []byte("PORT=8080\n"), 0644); err != nil {
		t.Fatal(err)
	}

	base := EnvMap{"HOST": "localhost", "PORT": "5432"}

	w := NewWatcher(20*time.Millisecond, p)
	w.Start()
	defer w.Stop()

	time.Sleep(40 * time.Millisecond)

	if err := os.WriteFile(p, []byte("PORT=9090\nDEBUG=true\n"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case ev := <-w.Events:
		merged := Merge(base, ev.Env)
		if merged["PORT"] != "9090" {
			t.Errorf("expected PORT=9090, got %s", merged["PORT"])
		}
		if merged["HOST"] != "localhost" {
			t.Errorf("expected HOST=localhost, got %s", merged["HOST"])
		}
		if merged["DEBUG"] != "true" {
			t.Errorf("expected DEBUG=true, got %s", merged["DEBUG"])
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout")
	}
}

package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeWatchEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeWatchEnv: %v", err)
	}
	return p
}

func TestWatcher_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	p := writeWatchEnv(t, dir, ".env", "KEY=original\n")

	w := NewWatcher(20*time.Millisecond, p)
	w.Start()
	defer w.Stop()

	// seed initial hash
	time.Sleep(40 * time.Millisecond)

	if err := os.WriteFile(p, []byte("KEY=changed\n"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case ev := <-w.Events:
		if ev.Path != p {
			t.Errorf("expected path %s, got %s", p, ev.Path)
		}
		if ev.Env["KEY"] != "changed" {
			t.Errorf("expected KEY=changed, got %s", ev.Env["KEY"])
		}
		if ev.OldHash == ev.NewHash {
			t.Error("expected hashes to differ")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for watch event")
	}
}

func TestWatcher_NoEventOnNoChange(t *testing.T) {
	dir := t.TempDir()
	p := writeWatchEnv(t, dir, ".env", "STABLE=1\n")

	w := NewWatcher(20*time.Millisecond, p)
	w.Start()
	defer w.Stop()

	time.Sleep(100 * time.Millisecond)

	select {
	case ev := <-w.Events:
		t.Errorf("unexpected event for unchanged file: %+v", ev)
	default:
		// correct — no event
	}
}

func TestWatcher_ErrorOnMissingFile(t *testing.T) {
	w := NewWatcher(20*time.Millisecond, "/nonexistent/.env")
	w.Start()
	defer w.Stop()

	select {
	case err := <-w.Errors:
		if err == nil {
			t.Error("expected non-nil error")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timeout waiting for error")
	}
}

func TestWatcher_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	p1 := writeWatchEnv(t, dir, "base.env", "A=1\n")
	p2 := writeWatchEnv(t, dir, "override.env", "B=2\n")

	w := NewWatcher(20*time.Millisecond, p1, p2)
	w.Start()
	defer w.Stop()

	time.Sleep(40 * time.Millisecond)

	if err := os.WriteFile(p2, []byte("B=updated\n"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case ev := <-w.Events:
		if ev.Path != p2 {
			t.Errorf("expected p2 to change, got %s", ev.Path)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout")
	}
}

package envfile

import (
	"testing"
)

func TestSnapshotStore_TakeAndGet(t *testing.T) {
	store := NewSnapshotStore()
	env := EnvMap{"FOO": "bar", "BAZ": "qux"}
	store.Take("v1", env)

	snap, err := store.Get("v1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if snap.Label != "v1" {
		t.Errorf("expected label v1, got %s", snap.Label)
	}
	if snap.Env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", snap.Env["FOO"])
	}
}

func TestSnapshotStore_GetMissing(t *testing.T) {
	store := NewSnapshotStore()
	_, err := store.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestSnapshotStore_IsolatesCopy(t *testing.T) {
	store := NewSnapshotStore()
	env := EnvMap{"KEY": "original"}
	store.Take("snap", env)
	env["KEY"] = "mutated"

	snap, _ := store.Get("snap")
	if snap.Env["KEY"] != "original" {
		t.Errorf("snapshot should be isolated from mutations, got %s", snap.Env["KEY"])
	}
}

func TestSnapshotStore_List(t *testing.T) {
	store := NewSnapshotStore()
	store.Take("a", EnvMap{"X": "1"})
	store.Take("b", EnvMap{"X": "2"})

	list := store.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(list))
	}
	if list[0].Label != "a" || list[1].Label != "b" {
		t.Errorf("unexpected order: %v", list)
	}
}

func TestSnapshotStore_Diff(t *testing.T) {
	store := NewSnapshotStore()
	store.Take("before", EnvMap{"FOO": "1", "BAR": "keep"})
	store.Take("after", EnvMap{"FOO": "2", "BAZ": "new"})

	entries, err := store.Diff("before", "after")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[string]DiffStatus{}
	for _, e := range entries {
		statuses[e.Key] = e.Status
	}

	if statuses["FOO"] != StatusChanged {
		t.Errorf("expected FOO changed, got %v", statuses["FOO"])
	}
	if statuses["BAR"] != StatusRemoved {
		t.Errorf("expected BAR removed, got %v", statuses["BAR"])
	}
	if statuses["BAZ"] != StatusAdded {
		t.Errorf("expected BAZ added, got %v", statuses["BAZ"])
	}
}

func TestSnapshot_Summary(t *testing.T) {
	store := NewSnapshotStore()
	snap := store.Take("test", EnvMap{"ALPHA": "1"})
	summary := snap.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	if len(summary) < 10 {
		t.Errorf("summary seems too short: %q", summary)
	}
}

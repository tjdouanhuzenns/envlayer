package envfile

import (
	"testing"
)

func TestDiff_Added(t *testing.T) {
	base := EnvMap{"A": "1"}
	next := EnvMap{"A": "1", "B": "2"}
	d := Diff(base, next)
	if len(d.Added) != 1 || d.Added["B"] != "2" {
		t.Errorf("expected B=2 in Added, got %v", d.Added)
	}
	if len(d.Removed) != 0 || len(d.Changed) != 0 {
		t.Errorf("unexpected changes: removed=%v changed=%v", d.Removed, d.Changed)
	}
}

func TestDiff_Removed(t *testing.T) {
	base := EnvMap{"A": "1", "B": "2"}
	next := EnvMap{"A": "1"}
	d := Diff(base, next)
	if len(d.Removed) != 1 || d.Removed["B"] != "2" {
		t.Errorf("expected B=2 in Removed, got %v", d.Removed)
	}
}

func TestDiff_Changed(t *testing.T) {
	base := EnvMap{"A": "old"}
	next := EnvMap{"A": "new"}
	d := Diff(base, next)
	if len(d.Changed) != 1 || d.Changed["A"] != "new" {
		t.Errorf("expected A=new in Changed, got %v", d.Changed)
	}
}

func TestDiff_Unchanged(t *testing.T) {
	base := EnvMap{"A": "1", "B": "2"}
	next := EnvMap{"A": "1", "B": "2"}
	d := Diff(base, next)
	if len(d.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged, got %v", d.Unchanged)
	}
	if d.HasChanges() {
		t.Error("expected HasChanges to be false")
	}
}

func TestDiff_NilBase(t *testing.T) {
	d := Diff(nil, EnvMap{"X": "1"})
	if len(d.Added) != 1 || d.Added["X"] != "1" {
		t.Errorf("expected X=1 in Added, got %v", d.Added)
	}
}

func TestDiff_NilNext(t *testing.T) {
	d := Diff(EnvMap{"X": "1"}, nil)
	if len(d.Removed) != 1 || d.Removed["X"] != "1" {
		t.Errorf("expected X=1 in Removed, got %v", d.Removed)
	}
}

func TestDiff_HasChanges(t *testing.T) {
	base := EnvMap{"A": "1"}
	next := EnvMap{"A": "2"}
	d := Diff(base, next)
	if !d.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

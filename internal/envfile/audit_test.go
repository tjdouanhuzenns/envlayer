package envfile

import (
	"strings"
	"testing"
)

func TestAuditMerge_Added(t *testing.T) {
	base := EnvMap{"A": "1"}
	layer := EnvMap{"A": "1", "B": "2"}
	log := NewAuditLog("test")
	result := AuditMerge(base, layer, log)
	if result["B"] != "2" {
		t.Errorf("expected B=2, got %s", result["B"])
	}
	if len(log.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(log.Entries))
	}
	var added *AuditEntry
	for i := range log.Entries {
		if log.Entries[i].Key == "B" {
			added = &log.Entries[i]
		}
	}
	if added == nil || added.Action != "added" {
		t.Errorf("expected B to be recorded as added")
	}
}

func TestAuditMerge_Changed(t *testing.T) {
	base := EnvMap{"HOST": "localhost"}
	layer := EnvMap{"HOST": "prod.example.com"}
	log := NewAuditLog("prod")
	AuditMerge(base, layer, log)
	if log.Entries[0].Action != "changed" {
		t.Errorf("expected changed, got %s", log.Entries[0].Action)
	}
	if log.Entries[0].OldValue != "localhost" {
		t.Errorf("unexpected old value: %s", log.Entries[0].OldValue)
	}
}

func TestAuditMerge_Unchanged(t *testing.T) {
	base := EnvMap{"PORT": "8080"}
	layer := EnvMap{"PORT": "8080"}
	log := NewAuditLog("staging")
	AuditMerge(base, layer, log)
	if log.Entries[0].Action != "unchanged" {
		t.Errorf("expected unchanged, got %s", log.Entries[0].Action)
	}
}

func TestAuditMerge_BasePreserved(t *testing.T) {
	base := EnvMap{"A": "1", "B": "2"}
	layer := EnvMap{"B": "99"}
	log := NewAuditLog("test")
	result := AuditMerge(base, layer, log)
	if result["A"] != "1" {
		t.Errorf("expected A=1 to be preserved")
	}
}

func TestAuditLog_Summary(t *testing.T) {
	base := EnvMap{"A": "old", "C": "same"}
	layer := EnvMap{"A": "new", "B": "added", "C": "same"}
	log := NewAuditLog("summary-test")
	AuditMerge(base, layer, log)
	summary := log.Summary()
	if !strings.Contains(summary, "summary-test") {
		t.Error("summary missing source label")
	}
	if !strings.Contains(summary, "+ B") {
		t.Error("summary missing added entry for B")
	}
	if !strings.Contains(summary, "~ A") {
		t.Error("summary missing changed entry for A")
	}
	if !strings.Contains(summary, "= C") {
		t.Error("summary missing unchanged entry for C")
	}
}

func TestNewAuditLog_Empty(t *testing.T) {
	log := NewAuditLog("empty")
	if len(log.Entries) != 0 {
		t.Error("new log should have no entries")
	}
	if log.Source != "empty" {
		t.Errorf("unexpected source: %s", log.Source)
	}
}

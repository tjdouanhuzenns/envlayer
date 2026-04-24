package envfile

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// AuditEntry records a single change made during a merge or layer operation.
type AuditEntry struct {
	Key      string
	OldValue string
	NewValue string
	Action   string // "added", "removed", "changed", "unchanged"
	Source   string
	At       time.Time
}

// AuditLog holds a collection of audit entries for a session.
type AuditLog struct {
	Entries []AuditEntry
	Source  string
}

// NewAuditLog creates an AuditLog tagged with the given source label.
func NewAuditLog(source string) *AuditLog {
	return &AuditLog{Source: source}
}

// Record appends a new entry to the audit log.
func (a *AuditLog) Record(key, oldVal, newVal, action string) {
	a.Entries = append(a.Entries, AuditEntry{
		Key:      key,
		OldValue: oldVal,
		NewValue: newVal,
		Action:   action,
		Source:   a.Source,
		At:       time.Now(),
	})
}

// AuditMerge compares base and layer, records changes, and returns the merged map.
func AuditMerge(base, layer EnvMap, log *AuditLog) EnvMap {
	result := make(EnvMap)
	for k, v := range base {
		result[k] = v
	}
	for k, newVal := range layer {
		if oldVal, exists := base[k]; exists {
			if oldVal != newVal {
				log.Record(k, oldVal, newVal, "changed")
			} else {
				log.Record(k, oldVal, newVal, "unchanged")
			}
		} else {
			log.Record(k, "", newVal, "added")
		}
		result[k] = newVal
	}
	return result
}

// Summary returns a human-readable summary of the audit log.
func (a *AuditLog) Summary() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Audit log [source: %s] — %d entries\n", a.Source, len(a.Entries)))
	sorted := make([]AuditEntry, len(a.Entries))
	copy(sorted, a.Entries)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })
	for _, e := range sorted {
		switch e.Action {
		case "added":
			sb.WriteString(fmt.Sprintf("  + %s = %q\n", e.Key, e.NewValue))
		case "removed":
			sb.WriteString(fmt.Sprintf("  - %s (was %q)\n", e.Key, e.OldValue))
		case "changed":
			sb.WriteString(fmt.Sprintf("  ~ %s: %q -> %q\n", e.Key, e.OldValue, e.NewValue))
		case "unchanged":
			sb.WriteString(fmt.Sprintf("  = %s\n", e.Key))
		}
	}
	return sb.String()
}

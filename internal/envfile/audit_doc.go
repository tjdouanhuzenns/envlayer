// Package envfile provides utilities for parsing, merging, diffing,
// validating, and exporting environment variable files.
//
// # Audit
//
// The audit module tracks changes that occur when merging environment maps.
// It records each key's old value, new value, and the action taken
// (added, removed, changed, or unchanged).
//
// Usage:
//
//	log := envfile.NewAuditLog("prod-overlay")
//	merged := envfile.AuditMerge(base, layer, log)
//	fmt.Print(log.Summary())
//
// Output format:
//
//	Audit log [source: prod-overlay] — 3 entries
//	  + NEW_KEY = "value"
//	  ~ CHANGED_KEY: "old" -> "new"
//	  = UNCHANGED_KEY
package envfile

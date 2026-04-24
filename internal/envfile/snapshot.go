package envfile

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Snapshot captures the state of an EnvMap at a point in time,
// along with metadata describing when and why it was taken.
type Snapshot struct {
	Label     string
	Timestamp time.Time
	Env       EnvMap
}

// SnapshotStore holds an ordered list of snapshots.
type SnapshotStore struct {
	snapshots []Snapshot
}

// NewSnapshotStore returns an empty SnapshotStore.
func NewSnapshotStore() *SnapshotStore {
	return &SnapshotStore{}
}

// Take records a deep copy of env under the given label.
func (s *SnapshotStore) Take(label string, env EnvMap) Snapshot {
	copy := make(EnvMap, len(env))
	for k, v := range env {
		copy[k] = v
	}
	snap := Snapshot{
		Label:     label,
		Timestamp: time.Now(),
		Env:       copy,
	}
	s.snapshots = append(s.snapshots, snap)
	return snap
}

// List returns all snapshots in order of creation.
func (s *SnapshotStore) List() []Snapshot {
	return s.snapshots
}

// Get returns the snapshot with the given label, or an error if not found.
func (s *SnapshotStore) Get(label string) (Snapshot, error) {
	for _, snap := range s.snapshots {
		if snap.Label == label {
			return snap, nil
		}
	}
	return Snapshot{}, fmt.Errorf("snapshot %q not found", label)
}

// Diff compares two snapshots and returns the DiffResult between them.
func (s *SnapshotStore) Diff(fromLabel, toLabel string) ([]DiffEntry, error) {
	from, err := s.Get(fromLabel)
	if err != nil {
		return nil, fmt.Errorf("from: %w", err)
	}
	to, err := s.Get(toLabel)
	if err != nil {
		return nil, fmt.Errorf("to: %w", err)
	}
	return Diff(from.Env, to.Env), nil
}

// Summary returns a human-readable summary of a snapshot.
func (snap Snapshot) Summary() string {
	keys := make([]string, 0, len(snap.Env))
	for k := range snap.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	fmt.Fprintf(&sb, "[%s] %s (%d keys)\n", snap.Timestamp.Format(time.RFC3339), snap.Label, len(keys))
	for _, k := range keys {
		fmt.Fprintf(&sb, "  %s=%s\n", k, snap.Env[k])
	}
	return sb.String()
}

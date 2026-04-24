// Package envfile provides snapshot functionality for capturing and comparing
// environment variable states over time.
//
// # Snapshot
//
// A Snapshot records a labelled, timestamped copy of an EnvMap. Snapshots are
// stored in a SnapshotStore and can be diffed against each other to understand
// how an environment changed between two points — for example, before and after
// a merge operation.
//
// Basic usage:
//
//	store := envfile.NewSnapshotStore()
//	store.Take("before", baseEnv)
//	merged := envfile.Merge(baseEnv, layerEnv)
//	store.Take("after", merged)
//
//	entries, err := store.Diff("before", "after")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, e := range entries {
//		fmt.Println(e.Key, e.Status)
//	}
package envfile

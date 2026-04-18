package envfile

// DiffResult holds the differences between two EnvMaps.
type DiffResult struct {
	Added    map[string]string // keys present in next but not in base
	Removed  map[string]string // keys present in base but not in next
	Changed  map[string]string // keys present in both but with different values (value is next value)
	Unchanged map[string]string // keys present in both with same value
}

// Diff computes the difference between a base EnvMap and a next EnvMap.
// base represents the original environment, next represents the updated one.
func Diff(base, next EnvMap) DiffResult {
	result := DiffResult{
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Changed:   make(map[string]string),
		Unchanged: make(map[string]string),
	}

	if base == nil {
		base = EnvMap{}
	}
	if next == nil {
		next = EnvMap{}
	}

	for k, v := range next {
		if baseVal, ok := base[k]; !ok {
			result.Added[k] = v
		} else if baseVal != v {
			result.Changed[k] = v
		} else {
			result.Unchanged[k] = v
		}
	}

	for k, v := range base {
		if _, ok := next[k]; !ok {
			result.Removed[k] = v
		}
	}

	return result
}

// HasChanges returns true if the DiffResult contains any added, removed, or changed keys.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

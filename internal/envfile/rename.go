package envfile

import (
	"fmt"
	"strings"
)

// RenameOptions controls how keys are renamed in an EnvMap.
type RenameOptions struct {
	// Mapping is a map of oldKey -> newKey.
	Mapping map[string]string

	// FailOnMissing causes Rename to return an error if a source key
	// from the mapping does not exist in the input map.
	FailOnMissing bool

	// DropOriginal removes the old key after renaming. Defaults to true
	// when not explicitly set via the functional options.
	DropOriginal bool
}

// DefaultRenameOptions returns a RenameOptions with sensible defaults.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{
		DropOriginal: true,
	}
}

// Rename returns a new EnvMap with keys renamed according to opts.Mapping.
// Keys not present in the mapping are carried over unchanged.
// If DropOriginal is true (the default), the old key is removed from the result.
func Rename(env EnvMap, opts RenameOptions) (EnvMap, error) {
	if env == nil {
		return nil, fmt.Errorf("rename: input map must not be nil")
	}
	if len(opts.Mapping) == 0 {
		return nil, fmt.Errorf("rename: mapping must not be empty")
	}

	// Validate mapping values are non-empty.
	for old, newKey := range opts.Mapping {
		if strings.TrimSpace(newKey) == "" {
			return nil, fmt.Errorf("rename: target key for %q must not be empty", old)
		}
	}

	result := make(EnvMap, len(env))
	for k, v := range env {
		result[k] = v
	}

	for oldKey, newKey := range opts.Mapping {
		val, exists := result[oldKey]
		if !exists {
			if opts.FailOnMissing {
				return nil, fmt.Errorf("rename: key %q not found in map", oldKey)
			}
			continue
		}
		result[newKey] = val
		if opts.DropOriginal && oldKey != newKey {
			delete(result, oldKey)
		}
	}

	return result, nil
}

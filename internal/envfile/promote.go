package envfile

import "fmt"

// PromoteOptions configures how promotion between environments is performed.
type PromoteOptions struct {
	// Keys to include in the promotion (empty means all keys)
	IncludeKeys []string
	// Keys to exclude from the promotion
	ExcludeKeys []string
	// If true, keys in the target that are not in the source are preserved
	PreserveTarget bool
}

// PromoteResult holds the outcome of a promotion operation.
type PromoteResult struct {
	Merged  EnvMap
	Added   []string
	Updated []string
	Skipped []string
}

// Promote merges selected keys from src into dst, respecting PromoteOptions.
// It returns a PromoteResult describing what changed.
func Promote(src, dst EnvMap, opts PromoteOptions) (PromoteResult, error) {
	if src == nil {
		return PromoteResult{}, fmt.Errorf("promote: source env map must not be nil")
	}

	include := make(map[string]bool, len(opts.IncludeKeys))
	for _, k := range opts.IncludeKeys {
		include[k] = true
	}
	exclude := make(map[string]bool, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		exclude[k] = true
	}

	result := NewEnvMap()
	if opts.PreserveTarget && dst != nil {
		for k, v := range dst {
			result[k] = v
		}
	}

	var added, updated, skipped []string

	for k, v := range src {
		if exclude[k] {
			skipped = append(skipped, k)
			continue
		}
		if len(include) > 0 && !include[k] {
			skipped = append(skipped, k)
			continue
		}
		_, exists := result[k]
		if exists {
			updated = append(updated, k)
		} else {
			added = append(added, k)
		}
		result[k] = v
	}

	return PromoteResult{
		Merged:  result,
		Added:   added,
		Updated: updated,
		Skipped: skipped,
	}, nil
}

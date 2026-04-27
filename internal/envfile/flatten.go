package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenOptions controls how nested key structures are flattened.
type FlattenOptions struct {
	// Separator is the string used to join key segments (default: "_").
	Separator string
	// UpperCase converts all resulting keys to uppercase.
	UpperCase bool
	// Prefix is an optional prefix prepended to every key.
	Prefix string
}

// DefaultFlattenOptions returns sensible defaults for FlattenOptions.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator: "_",
		UpperCase: true,
	}
}

// Flatten takes an EnvMap whose keys may contain a hierarchical separator
// (e.g. "DB.HOST" or "DB__HOST") and rewrites them using the target
// separator defined in opts, optionally uppercasing and prefixing them.
//
// The source separator is auto-detected: "." is tried first, then "__".
// Keys that contain neither are passed through unchanged (after case/prefix
// transforms).
func Flatten(env EnvMap, opts FlattenOptions) (EnvMap, error) {
	if env == nil {
		return nil, fmt.Errorf("flatten: input env map is nil")
	}

	sep := opts.Separator
	if sep == "" {
		sep = "_"
	}

	result := NewEnvMap()

	// Collect keys for deterministic ordering.
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]

		newKey := rewriteKey(k, sep)

		if opts.UpperCase {
			newKey = strings.ToUpper(newKey)
		}
		if opts.Prefix != "" {
			newKey = opts.Prefix + newKey
		}

		if _, exists := result[newKey]; exists {
			return nil, fmt.Errorf("flatten: key collision after rewrite: %q", newKey)
		}
		result[newKey] = v
	}

	return result, nil
}

// rewriteKey replaces recognised hierarchy separators with targetSep.
func rewriteKey(key, targetSep string) string {
	switch {
	case strings.Contains(key, "."):
		return strings.ReplaceAll(key, ".", targetSep)
	case strings.Contains(key, "__"):
		return strings.ReplaceAll(key, "__", targetSep)
	default:
		return key
	}
}

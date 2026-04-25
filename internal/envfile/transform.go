package envfile

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms an environment variable value.
type TransformFunc func(key, value string) (string, error)

// TransformOptions controls how transformations are applied.
type TransformOptions struct {
	// Keys limits transformation to specific keys; empty means all keys.
	Keys []string
	// FailOnError causes Transform to return early on the first error.
	FailOnError bool
}

// TransformResult holds the result of a single key transformation.
type TransformResult struct {
	Key      string
	OldValue string
	NewValue string
	Err      error
}

// Transform applies fn to each value in env, returning a new EnvMap and a
// slice of results describing what changed or failed.
func Transform(env EnvMap, fn TransformFunc, opts TransformOptions) (EnvMap, []TransformResult, error) {
	if env == nil {
		return NewEnvMap(), nil, nil
	}

	targetKeys := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		targetKeys[k] = true
	}

	out := NewEnvMap()
	results := make([]TransformResult, 0, len(env))

	for k, v := range env {
		if len(targetKeys) > 0 && !targetKeys[k] {
			out[k] = v
			continue
		}

		newVal, err := fn(k, v)
		result := TransformResult{Key: k, OldValue: v, NewValue: newVal, Err: err}
		results = append(results, result)

		if err != nil {
			out[k] = v // preserve original on error
			if opts.FailOnError {
				return out, results, fmt.Errorf("transform failed on key %q: %w", k, err)
			}
			continue
		}
		out[k] = newVal
	}

	return out, results, nil
}

// BuiltinTransforms provides common ready-to-use TransformFunc values.
var BuiltinTransforms = struct {
	UpperCase  TransformFunc
	LowerCase  TransformFunc
	TrimSpace  TransformFunc
	PrefixKeys func(prefix string) TransformFunc
}{
	UpperCase: func(_, v string) (string, error) { return strings.ToUpper(v), nil },
	LowerCase: func(_, v string) (string, error) { return strings.ToLower(v), nil },
	TrimSpace: func(_, v string) (string, error) { return strings.TrimSpace(v), nil },
	PrefixKeys: func(prefix string) TransformFunc {
		return func(_, v string) (string, error) { return prefix + v, nil }
	},
}

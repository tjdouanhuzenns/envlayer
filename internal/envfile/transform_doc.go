// Package envfile provides utilities for parsing, merging, and transforming
// environment variable files.
//
// # Transform
//
// Transform applies a user-supplied TransformFunc to every value (or a subset
// of keys) in an EnvMap, returning a new EnvMap with the transformed values
// and a slice of TransformResult entries describing each change.
//
// Built-in helpers are available via BuiltinTransforms:
//
//	env, _, _ := envfile.Transform(base, envfile.BuiltinTransforms.TrimSpace, envfile.TransformOptions{})
//
// Custom transforms can be composed:
//
//	upper := func(k, v string) (string, error) {
//	    return strings.ToUpper(v), nil
//	}
//	env, results, err := envfile.Transform(base, upper, envfile.TransformOptions{FailOnError: true})
package envfile

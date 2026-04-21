// Package envfile provides utilities for parsing, merging, diffing,
// validating, exporting, interpolating, and masking environment variable files.
//
// # Masking
//
// The Mask function returns a copy of an EnvMap with sensitive values
// replaced by a mask string. Sensitivity is determined by:
//
//   - Exact key matches via MaskOptions.Keys
//   - Substring pattern matches via MaskOptions.Patterns (case-insensitive)
//
// If no MaskOptions are provided, a built-in set of common patterns is used
// (e.g. SECRET, PASSWORD, TOKEN, API_KEY).
//
// Example:
//
//	masked := envfile.Mask(env, nil)
//	// env["DB_PASSWORD"] => "******"
//	// env["APP_NAME"]    => unchanged
//
// To reveal the last N characters of a masked value, set VisibleChars:
//
//	masked := envfile.Mask(env, &envfile.MaskOptions{
//		VisibleChars: 4,
//	})
//	// "mysupersecrettoken" => "******oken"
package envfile

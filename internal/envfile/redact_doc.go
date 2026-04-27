// Package envfile provides utilities for parsing, merging, and managing
// environment variable files across multiple deployment contexts.
//
// # Redact
//
// The Redact function produces a sanitised copy of an EnvMap by replacing
// sensitive values with a configurable placeholder string.
//
// Keys can be targeted in two ways:
//
//  1. Exact match — supply key names in RedactOptions.Keys. Matching is
//     case-insensitive so "db_password" and "DB_PASSWORD" are equivalent.
//
//  2. Pattern match — supply Go regular expressions in RedactOptions.Patterns.
//     Any key whose name matches at least one pattern will be redacted.
//
// The placeholder defaults to "[REDACTED]" when RedactOptions.Placeholder is
// left empty.
//
// Redact never modifies the original EnvMap; it always returns a new map.
//
// Example:
//
//	out, err := envfile.Redact(env, envfile.RedactOptions{
//		Patterns:    []string{"(?i)secret|password|token"},
//		Placeholder: "***",
//	})
package envfile

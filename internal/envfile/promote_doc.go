// Package envfile provides the Promote function for safely advancing
// environment variable sets from one deployment stage to another
// (e.g. staging → production).
//
// # Overview
//
// Promote copies keys from a source EnvMap into a destination EnvMap,
// with fine-grained control over which keys are included or excluded.
// It returns a PromoteResult that describes exactly which keys were
// added, updated, or skipped, making it easy to audit promotions.
//
// # Usage
//
//	result, err := envfile.Promote(stagingEnv, prodEnv, envfile.PromoteOptions{
//		ExcludeKeys:    []string{"DEBUG", "DEV_ONLY_FLAG"},
//		PreserveTarget: true,
//	})
//
// # Options
//
//   - IncludeKeys: promote only the listed keys (empty = all keys).
//   - ExcludeKeys: never promote the listed keys.
//   - PreserveTarget: keep destination keys that are absent in source.
package envfile

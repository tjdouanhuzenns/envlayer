// Package envfile provides the Compare function for side-by-side
// analysis of two environment maps.
//
// Compare categorises every key across two named EnvMaps into four
// buckets:
//
//   - OnlyLeft  – keys present only in the left (e.g. dev) map
//   - OnlyRight – keys present only in the right (e.g. prod) map
//   - Differ    – keys present in both but with different values
//   - Same      – keys present in both with identical values
//
// Example:
//
//	dev  := envfile.EnvMap{"HOST": "localhost", "DEBUG": "true"}
//	prod := envfile.EnvMap{"HOST": "prod.example.com"}
//
//	result := envfile.Compare("dev", dev, "prod", prod)
//	fmt.Print(result.Summary())
package envfile

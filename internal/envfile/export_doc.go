// Package envfile provides utilities for parsing, merging, diffing,
// validating, resolving, and exporting environment variable files.
//
// # Export
//
// The Export feature allows serializing a merged EnvMap into multiple
// output formats suitable for different tooling needs:
//
//   - FormatDotenv: standard KEY=VALUE format, compatible with most tools
//   - FormatExport: shell-sourceable `export KEY=VALUE` format
//   - FormatJSON:   JSON object format for use with config systems
//
// Example usage:
//
//	env, err := envfile.ResolveAndMerge(".", "prod")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := envfile.Export(os.Stdout, env, envfile.FormatExport); err != nil {
//		log.Fatal(err)
//	}
package envfile

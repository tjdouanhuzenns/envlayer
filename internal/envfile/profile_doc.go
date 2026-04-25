// Package envfile provides the Profile and ProfileRegistry types for
// managing named, layered environment configurations.
//
// A Profile is a named sequence of .env file paths that are merged in
// order — later layers override earlier ones. Profiles are stored in a
// ProfileRegistry which can resolve any registered profile into a final
// merged EnvMap ready for use or export.
//
// Example:
//
//	reg := envfile.NewProfileRegistry()
//	reg.Register(&envfile.Profile{
//		Name:   "staging",
//		Layers: []string{".env", ".env.staging"},
//	})
//	env, err := reg.Resolve("staging")
package envfile

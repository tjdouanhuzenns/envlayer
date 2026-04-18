// Package envfile provides utilities for parsing, merging, and writing
// environment variable files (.env format).
//
// # Writing
//
// Use WriteFile to persist an EnvMap to disk, or WriteString to serialize
// it to a string. Both functions sort keys alphabetically for deterministic
// output and quote values that contain spaces, tabs, or special characters.
//
// Example:
//
//	env := envfile.EnvMap{"APP_ENV": "production", "PORT": "8080"}
//
//	// Write to file
//	if err := envfile.WriteFile(".env.out", env); err != nil {
//		log.Fatal(err)
//	}
//
//	// Or get as string
//	fmt.Print(envfile.WriteString(env))
package envfile

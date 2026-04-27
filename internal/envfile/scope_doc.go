// Package envfile provides scope-based environment layering.
//
// A ScopeRegistry allows you to register named environments (scopes) such as
// "base", "dev", "staging", or "prod", and then resolve them in a specified
// order — with later scopes overriding earlier ones.
//
// Example:
//
//	reg := envfile.NewScopeRegistry()
//	reg.Register("base", baseEnv)
//	reg.Register("prod", prodEnv)
//
//	merged, err := reg.Resolve("base", "prod")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// This is useful for building layered configuration pipelines where each
// environment inherits from a common base.
package envfile

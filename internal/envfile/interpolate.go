package envfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// varPattern matches ${VAR_NAME} and $VAR_NAME style references.
var varPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateOptions controls how interpolation resolves variables.
type InterpolateOptions struct {
	// FallbackToOS allows falling back to OS environment variables
	// when a key is not found in the provided map.
	FallbackToOS bool
	// FailOnMissing returns an error if a referenced variable cannot be resolved.
	FailOnMissing bool
}

// Interpolate resolves variable references within the values of env.
// References to other keys in env are substituted; if FallbackToOS is set,
// unresolved references are looked up in os.Getenv.
func Interpolate(env EnvMap, opts InterpolateOptions) (EnvMap, error) {
	result := NewEnvMap()
	for k, v := range env {
		resolved, err := interpolateValue(v, env, opts)
		if err != nil {
			return nil, fmt.Errorf("interpolating key %q: %w", k, err)
		}
		result[k] = resolved
	}
	return result, nil
}

func interpolateValue(value string, env EnvMap, opts InterpolateOptions) (string, error) {
	var firstErr error
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		if firstErr != nil {
			return match
		}
		// Extract the variable name from either capture group.
		subs := varPattern.FindStringSubmatch(match)
		name := subs[1]
		if name == "" {
			name = subs[2]
		}
		if val, ok := env[name]; ok {
			return val
		}
		if opts.FallbackToOS {
			if val, ok := os.LookupEnv(name); ok {
				return val
			}
		}
		if opts.FailOnMissing {
			firstErr = fmt.Errorf("variable %q is not defined", name)
			return match
		}
		// Leave unresolved references as empty string.
		return ""
	})
	if firstErr != nil {
		return "", firstErr
	}
	_ = strings // imported for potential future use
	return result, nil
}

package envfile

import (
	"fmt"
	"os"\n	"sort"
	"strings"
)

// WriteFile writes an EnvMap to a file in KEY=VALUE format.
// Values containing spaces or special characters are quoted.
func WriteFile(path string, env EnvMap) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("envfile: create %s: %w", path, err)
	}
	defer f.Close()

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		if needsQuoting(v) {
			v = `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
		}
		if _, err := fmt.Fprintf(f, "%s=%s\n", k, v); err != nil {
			return fmt.Errorf("envfile: write key %s: %w", k, err)
		}
	}
	return nil
}

// WriteString serializes an EnvMap to a string in KEY=VALUE format.
func WriteString(env EnvMap) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		if needsQuoting(v) {
			v = `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
		}
		sb.WriteString(k + "=" + v + "\n")
	}
	return sb.String()
}

func needsQuoting(v string) bool {
	if v == "" {
		return false
	}
	for _, c := range v {
		if c == ' ' || c == '\t' || c == '#' || c == '$' || c == '\'' || c == '"' {
			return true
		}
	}
	return false
}

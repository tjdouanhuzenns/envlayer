package envfile

import (
	"fmt"
	"os"
	"path/filepath"
)

// Context represents a named environment context (e.g. dev, staging, prod).
type Context struct {
	Name    string
	BaseDir string
}

// ResolveFiles returns an ordered list of .env file paths for a given context.
// It looks for: .env, .env.<context>, .env.<context>.local (in that order).
func ResolveFiles(baseDir, context string) ([]string, error) {
	candidates := []string{
		filepath.Join(baseDir, ".env"),
		filepath.Join(baseDir, fmt.Sprintf(".env.%s", context)),
		filepath.Join(baseDir, fmt.Sprintf(".env.%s.local", context)),
	}

	var resolved []string
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			resolved = append(resolved, path)
		}
	}

	if len(resolved) == 0 {
		return nil, fmt.Errorf("no env files found for context %q in %q", context, baseDir)
	}

	return resolved, nil
}

// ResolveAndMerge resolves all env files for a context and merges them in order.
// Later files override earlier ones.
func ResolveAndMerge(baseDir, context string) (EnvMap, error) {
	files, err := ResolveFiles(baseDir, context)
	if err != nil {
		return nil, err
	}

	maps := make([]EnvMap, 0, len(files))
	for _, f := range files {
		em, err := ParseFile(f)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %w", f, err)
		}
		maps = append(maps, em)
	}

	return MergeFiles(maps...), nil
}

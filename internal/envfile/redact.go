package envfile

import (
	"regexp"
	"strings"
)

// RedactOptions controls how values are redacted in an EnvMap.
type RedactOptions struct {
	// Keys is an explicit list of keys to redact.
	Keys []string
	// Patterns is a list of regex patterns matched against key names.
	Patterns []string
	// Placeholder is the string used to replace redacted values.
	// Defaults to "[REDACTED]" if empty.
	Placeholder string
}

// Redact returns a new EnvMap with sensitive values replaced by a placeholder.
// It matches keys by exact name (case-insensitive) or by regex pattern.
// The original EnvMap is never modified.
func Redact(env EnvMap, opts RedactOptions) (EnvMap, error) {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "[REDACTED]"
	}

	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}

	exactKeys := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		exactKeys[strings.ToUpper(k)] = struct{}{}
	}

	result := make(EnvMap, len(env))
	for k, v := range env {
		if shouldRedact(k, exactKeys, compiled) {
			result[k] = placeholder
		} else {
			result[k] = v
		}
	}
	return result, nil
}

func shouldRedact(key string, exactKeys map[string]struct{}, patterns []*regexp.Regexp) bool {
	if _, ok := exactKeys[strings.ToUpper(key)]; ok {
		return true
	}
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

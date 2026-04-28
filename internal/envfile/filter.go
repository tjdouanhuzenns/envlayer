package envfile

import (
	"regexp"
	"strings"
)

// FilterOptions controls how keys are selected from an EnvMap.
type FilterOptions struct {
	// Prefix keeps only keys that start with the given prefix.
	Prefix string

	// Suffix keeps only keys that end with the given suffix.
	Suffix string

	// Pattern keeps only keys matching the given regular expression.
	Pattern string

	// Keys is an explicit allowlist of keys to keep.
	Keys []string

	// Invert reverses the filter — keeps keys that do NOT match.
	Invert bool
}

// Filter returns a new EnvMap containing only the entries that satisfy
// the provided FilterOptions. Multiple criteria are combined with AND
// logic (all specified criteria must match). When Invert is true the
// logic is flipped and only non-matching keys are kept.
//
// An empty FilterOptions (all zero values) returns a shallow copy of
// the original map unchanged.
func Filter(env EnvMap, opts FilterOptions) (EnvMap, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
	}

	allowSet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		allowSet[k] = struct{}{}
	}

	result := NewEnvMap()
	for k, v := range env {
		matched := matches(k, opts.Prefix, opts.Suffix, re, allowSet)
		if opts.Invert {
			matched = !matched
		}
		if matched {
			result[k] = v
		}
	}
	return result, nil
}

func matches(key, prefix, suffix string, re *regexp.Regexp, allowSet map[string]struct{}) bool {
	if prefix != "" && !strings.HasPrefix(key, prefix) {
		return false
	}
	if suffix != "" && !strings.HasSuffix(key, suffix) {
		return false
	}
	if re != nil && !re.MatchString(key) {
		return false
	}
	if len(allowSet) > 0 {
		if _, ok := allowSet[key]; !ok {
			return false
		}
	}
	return true
}

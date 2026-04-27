package envfile

// CloneOptions controls how an EnvMap is cloned.
type CloneOptions struct {
	// Prefix filters keys: only keys with this prefix are included.
	// An empty string means all keys are cloned.
	Prefix string

	// StripPrefix removes the Prefix from each key in the resulting map.
	StripPrefix bool

	// ExcludeKeys lists exact key names to omit from the clone.
	ExcludeKeys []string
}

// Clone returns a deep copy of src, applying any options provided.
// If opts is nil, all keys are copied without modification.
func Clone(src EnvMap, opts *CloneOptions) (EnvMap, error) {
	if src == nil {
		return EnvMap{}, nil
	}

	exclude := make(map[string]struct{})
	if opts != nil {
		for _, k := range opts.ExcludeKeys {
			exclude[k] = struct{}{}
		}
	}

	out := make(EnvMap, len(src))

	for k, v := range src {
		// Apply exclusion list.
		if _, skip := exclude[k]; skip {
			continue
		}

		outKey := k

		if opts != nil && opts.Prefix != "" {
			// Filter by prefix.
			if len(k) < len(opts.Prefix) || k[:len(opts.Prefix)] != opts.Prefix {
				continue
			}
			if opts.StripPrefix {
				outKey = k[len(opts.Prefix):]
				if outKey == "" {
					// Key was exactly the prefix — skip it.
					continue
				}
			}
		}

		out[outKey] = v
	}

	return out, nil
}

package envfile

import "fmt"

// MergeStrategy controls how values are combined during a merge.
type MergeStrategy int

const (
	// StrategyOverride replaces base values with layer values (default).
	StrategyOverride MergeStrategy = iota
	// StrategyKeepBase preserves base values; layer only fills missing keys.
	StrategyKeepBase
)

// Merge combines one or more EnvMaps into a base map.
// Layers are applied left-to-right, so later layers take precedence.
func Merge(base *EnvMap, layers []*EnvMap, strategy MergeStrategy) (*EnvMap, error) {
	if base == nil {
		return nil, fmt.Errorf("base EnvMap must not be nil")
	}

	result := NewEnvMap()

	// Seed result with base.
	for _, k := range base.Keys {
		result.Set(k, base.Values[k])
	}

	for _, layer := range layers {
		if layer == nil {
			continue
		}
		for _, k := range layer.Keys {
			switch strategy {
			case StrategyOverride:
				result.Set(k, layer.Values[k])
			case StrategyKeepBase:
				if _, exists := result.Values[k]; !exists {
					result.Set(k, layer.Values[k])
				}
			}
		}
	}

	return result, nil
}

// MergeFiles is a convenience wrapper that parses and merges files in order.
// The first path is treated as the base; subsequent paths are layers.
func MergeFiles(paths []string, strategy MergeStrategy) (*EnvMap, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("at least one file path required")
	}

	base, err := ParseFile(paths[0])
	if err != nil {
		return nil, fmt.Errorf("parsing base file: %w", err)
	}

	var layers []*EnvMap
	for _, p := range paths[1:] {
		env, err := ParseFile(p)
		if err != nil {
			return nil, fmt.Errorf("parsing layer %s: %w", p, err)
		}
		layers = append(layers, env)
	}

	return Merge(base, layers, strategy)
}

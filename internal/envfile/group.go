package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// Group holds a named collection of environment variable keys.
type Group struct {
	Name string
	Keys []string
}

// GroupRegistry manages named groups of environment variable keys,
// allowing bulk operations (export, mask, validate) on logical subsets.
type GroupRegistry struct {
	groups map[string]*Group
}

// NewGroupRegistry returns an initialised, empty GroupRegistry.
func NewGroupRegistry() *GroupRegistry {
	return &GroupRegistry{groups: make(map[string]*Group)}
}

// Register adds a new named group with the given keys.
// Returns an error if the name is empty, already registered, or no keys are provided.
func (r *GroupRegistry) Register(name string, keys []string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("group name must not be empty")
	}
	if len(keys) == 0 {
		return fmt.Errorf("group %q must contain at least one key", name)
	}
	if _, exists := r.groups[name]; exists {
		return fmt.Errorf("group %q is already registered", name)
	}
	copy := make([]string, len(keys))
	for i, k := range keys {
		copy[i] = strings.TrimSpace(k)
	}
	r.groups[name] = &Group{Name: name, Keys: copy}
	return nil
}

// Get returns the Group registered under name, or an error if not found.
func (r *GroupRegistry) Get(name string) (*Group, error) {
	g, ok := r.groups[name]
	if !ok {
		return nil, fmt.Errorf("group %q not found", name)
	}
	return g, nil
}

// List returns all registered group names in sorted order.
func (r *GroupRegistry) List() []string {
	names := make([]string, 0, len(r.groups))
	for n := range r.groups {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Extract returns a new EnvMap containing only the keys belonging to the
// named group that are present in src. Missing keys are silently skipped
// unless strict is true, in which case an error is returned.
func (r *GroupRegistry) Extract(name string, src EnvMap, strict bool) (EnvMap, error) {
	g, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	out := NewEnvMap()
	for _, k := range g.Keys {
		v, ok := src[k]
		if !ok {
			if strict {
				return nil, fmt.Errorf("group %q: key %q not found in source", name, k)
			}
			continue
		}
		out[k] = v
	}
	return out, nil
}

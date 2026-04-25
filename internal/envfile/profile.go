package envfile

import (
	"fmt"
	"sort"
)

// Profile represents a named environment configuration combining
// multiple layers in a defined order.
type Profile struct {
	Name   string
	Layers []string // ordered list of env file paths
}

// ProfileRegistry holds named profiles for an envlayer project.
type ProfileRegistry struct {
	profiles map[string]*Profile
}

// NewProfileRegistry creates an empty ProfileRegistry.
func NewProfileRegistry() *ProfileRegistry {
	return &ProfileRegistry{
		profiles: make(map[string]*Profile),
	}
}

// Register adds or replaces a profile in the registry.
func (r *ProfileRegistry) Register(p *Profile) error {
	if p == nil {
		return fmt.Errorf("profile must not be nil")
	}
	if p.Name == "" {
		return fmt.Errorf("profile name must not be empty")
	}
	r.profiles[p.Name] = p
	return nil
}

// Get retrieves a profile by name.
func (r *ProfileRegistry) Get(name string) (*Profile, bool) {
	p, ok := r.profiles[name]
	return p, ok
}

// List returns sorted profile names.
func (r *ProfileRegistry) List() []string {
	names := make([]string, 0, len(r.profiles))
	for k := range r.profiles {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// Resolve merges all layers of the named profile and returns the final EnvMap.
func (r *ProfileRegistry) Resolve(name string) (EnvMap, error) {
	p, ok := r.profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile %q not found", name)
	}
	if len(p.Layers) == 0 {
		return NewEnvMap(), nil
	}
	return ResolveAndMerge(p.Layers...)
}

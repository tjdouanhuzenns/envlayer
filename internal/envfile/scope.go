package envfile

import "fmt"

// Scope represents a named environment context (e.g. dev, staging, prod).
type Scope struct {
	Name string
	Env  EnvMap
}

// ScopeRegistry holds multiple named scopes and supports layered resolution.
type ScopeRegistry struct {
	scopes []Scope
}

// NewScopeRegistry creates an empty ScopeRegistry.
func NewScopeRegistry() *ScopeRegistry {
	return &ScopeRegistry{}
}

// Register adds a named scope to the registry.
// Returns an error if the name is empty or already registered.
func (r *ScopeRegistry) Register(name string, env EnvMap) error {
	if name == "" {
		return fmt.Errorf("scope name must not be empty")
	}
	for _, s := range r.scopes {
		if s.Name == name {
			return fmt.Errorf("scope %q already registered", name)
		}
	}
	r.scopes = append(r.scopes, Scope{Name: name, Env: env})
	return nil
}

// Get returns the EnvMap for the named scope.
func (r *ScopeRegistry) Get(name string) (EnvMap, bool) {
	for _, s := range r.scopes {
		if s.Name == name {
			return s.Env, true
		}
	}
	return nil, false
}

// List returns all registered scope names in registration order.
func (r *ScopeRegistry) List() []string {
	names := make([]string, len(r.scopes))
	for i, s := range r.scopes {
		names[i] = s.Name
	}
	return names
}

// Resolve merges scopes in order, with later scopes overriding earlier ones.
// Returns an error if any named scope is not found.
func (r *ScopeRegistry) Resolve(names ...string) (EnvMap, error) {
	result := NewEnvMap()
	for _, name := range names {
		env, ok := r.Get(name)
		if !ok {
			return nil, fmt.Errorf("scope %q not found", name)
		}
		result = Merge(result, env)
	}
	return result, nil
}

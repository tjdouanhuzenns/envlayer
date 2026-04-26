package envfile

import "fmt"

// PipelineStep represents a single transformation step in a pipeline.
type PipelineStep struct {
	Name string
	Fn   func(EnvMap) (EnvMap, error)
}

// PipelineResult holds the output of each step in a pipeline run.
type PipelineResult struct {
	// Steps contains the EnvMap state after each named step.
	Steps map[string]EnvMap
	// Final is the resulting EnvMap after all steps have been applied.
	Final EnvMap
}

// Pipeline is an ordered sequence of transformation steps applied to an EnvMap.
// Each step receives the output of the previous step, allowing composable
// processing chains such as merge → interpolate → validate → mask → export.
//
// Example:
//
//	p := Pipeline{
//		{Name: "interpolate", Fn: func(m EnvMap) (EnvMap, error) { return Interpolate(m, false) }},
//		{Name: "mask",        Fn: func(m EnvMap) (EnvMap, error) { return Mask(m, MaskOptions{}) }},
//	}
//	res, err := p.Run(base)
type Pipeline []PipelineStep

// Run executes each step in the pipeline sequentially, starting from src.
// It returns a PipelineResult containing the intermediate state after each
// named step as well as the final merged result.
//
// If any step returns an error the pipeline halts immediately and the error
// is wrapped with the failing step name for easy diagnosis.
func (p Pipeline) Run(src EnvMap) (PipelineResult, error) {
	result := PipelineResult{
		Steps: make(map[string]EnvMap, len(p)),
	}

	current := src.Clone()

	for _, step := range p {
		out, err := step.Fn(current)
		if err != nil {
			return result, fmt.Errorf("pipeline step %q: %w", step.Name, err)
		}
		result.Steps[step.Name] = out.Clone()
		current = out
	}

	result.Final = current
	return result, nil
}

// Clone returns a shallow copy of the EnvMap so that pipeline steps cannot
// accidentally mutate earlier snapshots stored in PipelineResult.Steps.
func (m EnvMap) Clone() EnvMap {
	copy := make(EnvMap, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return copy
}

// NewPipeline constructs a Pipeline from a variadic list of steps.
func NewPipeline(steps ...PipelineStep) Pipeline {
	return Pipeline(steps)
}

// Step is a convenience constructor for a PipelineStep.
func Step(name string, fn func(EnvMap) (EnvMap, error)) PipelineStep {
	return PipelineStep{Name: name, Fn: fn}
}

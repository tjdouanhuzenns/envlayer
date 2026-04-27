// Package envfile provides pipeline support for chaining environment
// variable transformations in a composable, step-by-step fashion.
//
// # Pipeline
//
// A Pipeline allows you to chain multiple EnvMap transformation steps
// together. Each step receives the output of the previous step, enabling
// complex workflows such as merge → interpolate → mask → export to be
// expressed cleanly.
//
// # Basic Usage
//
//	p := envfile.NewPipeline(baseEnv).
//		Step("interpolate", func(m envfile.EnvMap) (envfile.EnvMap, error) {
//			return envfile.Interpolate(m, envfile.InterpolateOptions{})
//		}).
//		Step("mask", func(m envfile.EnvMap) (envfile.EnvMap, error) {
//			return envfile.Mask(m, envfile.MaskOptions{})
//		})
//
//	result, err := p.Run()
//
// # Error Handling
//
// If any step returns an error, the pipeline halts immediately and
// returns the error along with the name of the failing step. The
// original input map is never mutated.
package envfile

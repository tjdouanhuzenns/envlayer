// Package envfile provides schema validation for environment variable maps.
//
// Schema validation allows you to define expected keys, whether they are
// required, and optional constraints such as regex patterns or allowed value
// lists. This is useful for ensuring that a merged or resolved EnvMap satisfies
// the requirements of a given deployment environment before it is used.
//
// Example usage:
//
//	schema := envfile.Schema{
//		Fields: []envfile.SchemaField{
//			{Key: "APP_ENV", Required: true, AllowedValues: []string{"dev", "staging", "prod"}},
//			{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
//			{Key: "LOG_LEVEL", Required: false},
//		},
//	}
//
//	errs := envfile.ValidateSchema(env, schema)
//	if len(errs) > 0 {
//		for _, e := range errs {
//			fmt.Println(e)
//		}
//	}
package envfile

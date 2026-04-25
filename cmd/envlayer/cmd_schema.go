package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"envlayer/internal/envfile"
)

var schemaCmd = &cobra.Command{
	Use:   "schema [env-file]",
	Short: "Validate an env file against a built-in schema definition",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := envfile.ParseFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to parse env file: %w", err)
		}

		// Example built-in schema; in a real tool this would be loaded from a
		// .envschema file or flags.
		schema := envfile.Schema{
			Fields: []envfile.SchemaField{
				{
					Key:           "APP_ENV",
					Required:      true,
					AllowedValues: []string{"dev", "staging", "prod"},
				},
				{
					Key:      "PORT",
					Required: false,
					Pattern:  regexp.MustCompile(`^\d+$`),
				},
			},
		}

		errs := envfile.ValidateSchema(env, schema)
		if len(errs) == 0 {
			fmt.Println("schema validation passed")
			return nil
		}

		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  - %s\n", e.Error())
		}
		return fmt.Errorf("schema validation failed with %d error(s)", len(errs))
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)
}

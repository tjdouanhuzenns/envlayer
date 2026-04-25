package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envlayer/internal/envfile"
)

func init() {
	var files []string
	var warnEmpty bool
	var warnCaps bool
	var warnUnderscore bool
	var strict bool

	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint environment files for style and quality issues",
		Example: `  envlayer lint --file .env --file .env.local
  envlayer lint --file .env.prod --strict`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(files) == 0 {
				return fmt.Errorf("at least one --file is required")
			}

			merged, err := envfile.MergeFiles(files...)
			if err != nil {
				return fmt.Errorf("failed to load env files: %w", err)
			}

			opts := envfile.LintOptions{
				WarnEmptyValues:       warnEmpty,
				WarnNotAllCaps:        warnCaps,
				WarnLeadingUnderscore: warnUnderscore,
			}

			issues := envfile.Lint(merged, opts)

			if len(issues) == 0 {
				fmt.Println("✔ No lint issues found.")
				return nil
			}

			for _, issue := range issues {
				fmt.Fprintln(os.Stderr, issue.String())
			}

			fmt.Fprintf(os.Stderr, "\n%d issue(s) found.\n", len(issues))

			if strict {
				return fmt.Errorf("lint failed with %d issue(s)", len(issues))
			}
			return nil
		},
	}

	lintCmd.Flags().StringArrayVarP(&files, "file", "f", nil, "env file(s) to lint (merged in order)")
	lintCmd.Flags().BoolVar(&warnEmpty, "warn-empty", true, "warn on empty values")
	lintCmd.Flags().BoolVar(&warnCaps, "warn-caps", true, "warn on keys that are not all uppercase")
	lintCmd.Flags().BoolVar(&warnUnderscore, "warn-underscore", false, "warn on keys with a leading underscore")
	lintCmd.Flags().BoolVar(&strict, "strict", false, "exit with non-zero status if any issues are found")

	rootCmd.AddCommand(lintCmd)
}

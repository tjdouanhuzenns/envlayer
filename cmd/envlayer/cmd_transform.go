package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envlayer/internal/envfile"
)

func init() {
	var (
		inputFile  string
		outputFile string
		operation  string
		keys       []string
		prefix     string
	)

	cmd := &cobra.Command{
		Use:   "transform",
		Short: "Apply a transformation to environment variable values",
		Long: `Transform applies a built-in operation to values in an env file.

Operations: uppercase, lowercase, trim, prefix`,
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.ParseFile(inputFile)
			if err != nil {
				return fmt.Errorf("parse %q: %w", inputFile, err)
			}

			var fn envfile.TransformFunc
			switch strings.ToLower(operation) {
			case "uppercase":
				fn = envfile.BuiltinTransforms.UpperCase
			case "lowercase":
				fn = envfile.BuiltinTransforms.LowerCase
			case "trim":
				fn = envfile.BuiltinTransforms.TrimSpace
			case "prefix":
				if prefix == "" {
					return fmt.Errorf("--prefix is required for the 'prefix' operation")
				}
				fn = envfile.BuiltinTransforms.PrefixKeys(prefix)
			default:
				return fmt.Errorf("unknown operation %q; choose uppercase, lowercase, trim, or prefix", operation)
			}

			opts := envfile.TransformOptions{Keys: keys, FailOnError: true}
			out, results, err := envfile.Transform(env, fn, opts)
			if err != nil {
				return err
			}

			for _, r := range results {
				if r.OldValue != r.NewValue {
					fmt.Fprintf(os.Stderr, "  ~ %s: %q -> %q\n", r.Key, r.OldValue, r.NewValue)
				}
			}

			if outputFile != "" {
				return envfile.WriteFile(outputFile, out)
			}
			s, err := envfile.WriteString(out)
			if err != nil {
				return err
			}
			fmt.Print(s)
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputFile, "file", "f", ".env", "Input env file")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	cmd.Flags().StringVarP(&operation, "op", "p", "trim", "Operation to apply (uppercase|lowercase|trim|prefix)")
	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "Limit transformation to these keys (comma-separated)")
	cmd.Flags().StringVar(&prefix, "prefix", "", "Prefix string for the 'prefix' operation")

	rootCmd.AddCommand(cmd)
}

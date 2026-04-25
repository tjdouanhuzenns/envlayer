package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/user/envlayer/internal/envfile"
)

func init() {
	var includeKeys []string
	var excludeKeys []string
	var preserveTarget bool
	var outputFile string

	promoteCmd := &cobra.Command{
		Use:   "promote <source-file> <target-file>",
		Short: "Promote env vars from one environment file into another",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcPath, dstPath := args[0], args[1]

			src, err := envfile.ParseFile(srcPath)
			if err != nil {
				return fmt.Errorf("reading source %q: %w", srcPath, err)
			}

			var dst envfile.EnvMap
			if _, statErr := os.Stat(dstPath); statErr == nil {
				dst, err = envfile.ParseFile(dstPath)
				if err != nil {
					return fmt.Errorf("reading target %q: %w", dstPath, err)
				}
			}

			result, err := envfile.Promote(src, dst, envfile.PromoteOptions{
				IncludeKeys:    includeKeys,
				ExcludeKeys:    excludeKeys,
				PreserveTarget: preserveTarget,
			})
			if err != nil {
				return err
			}

			out := dstPath
			if outputFile != "" {
				out = outputFile
			}
			if err := envfile.WriteFile(out, result.Merged); err != nil {
				return fmt.Errorf("writing output: %w", err)
			}

			fmt.Printf("Promoted to %s\n", out)
			if len(result.Added) > 0 {
				fmt.Printf("  Added:   %s\n", strings.Join(result.Added, ", "))
			}
			if len(result.Updated) > 0 {
				fmt.Printf("  Updated: %s\n", strings.Join(result.Updated, ", "))
			}
			if len(result.Skipped) > 0 {
				fmt.Printf("  Skipped: %s\n", strings.Join(result.Skipped, ", "))
			}
			return nil
		},
	}

	promoteCmd.Flags().StringSliceVar(&includeKeys, "include", nil, "keys to include (default: all)")
	promoteCmd.Flags().StringSliceVar(&excludeKeys, "exclude", nil, "keys to exclude")
	promoteCmd.Flags().BoolVar(&preserveTarget, "preserve-target", true, "keep target keys not present in source")
	promoteCmd.Flags().StringVarP(&outputFile, "output", "o", "", "write result to this file instead of target")

	rootCmd.AddCommand(promoteCmd)
}

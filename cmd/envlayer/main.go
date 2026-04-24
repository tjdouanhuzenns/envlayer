// Package main is the entry point for the envlayer CLI tool.
// It wires together the internal envfile packages and exposes
// a user-facing command-line interface for merging, diffing,
// exporting, validating, and inspecting environment variable files.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// rootCmd builds the top-level cobra command and registers all sub-commands.
func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envlayer",
		Short: "Layer and merge environment variable files across contexts",
		Long: `envlayer is a CLI tool for composing .env files across
development, staging, and production contexts.

It supports merging layers, diffing changes, exporting to multiple
formats, validating keys, masking secrets, and snapshotting state.`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(
		newMergeCmd(),
		newDiffCmd(),
		newExportCmd(),
		newValidateCmd(),
		newSnapshotCmd(),
	)

	return root
}

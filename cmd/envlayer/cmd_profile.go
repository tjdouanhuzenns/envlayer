package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envlayer/internal/envfile"
)

func init() {
	var layers []string
	var format string

	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage and resolve named environment profiles",
	}

	resolveCmd := &cobra.Command{
		Use:   "resolve <name>",
		Short: "Resolve a profile by merging its layers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reg := envfile.NewProfileRegistry()
			if err := reg.Register(&envfile.Profile{
				Name:   args[0],
				Layers: layers,
			}); err != nil {
				return err
			}
			env, err := reg.Resolve(args[0])
			if err != nil {
				return err
			}
			fmt_parsed, err := envfile.ParseFormat(format)
			if err != nil {
				return err
			}
			out, err := envfile.Export(env, fmt_parsed)
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		},
	}
	resolveCmd.Flags().StringSliceVarP(&layers, "layers", "l", nil, "Ordered list of env files (comma-separated or repeated)")
	resolveCmd.Flags().StringVarP(&format, "format", "f", "dotenv", "Output format: dotenv, export, json")
	_ = resolveCmd.MarkFlagRequired("layers")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered profiles from env files",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Discover profiles from conventional filenames in cwd
			entries, err := os.ReadDir(".")
			if err != nil {
				return err
			}
			for _, e := range entries {
				name := e.Name()
				if strings.HasPrefix(name, ".env.") {
					suffix := strings.TrimPrefix(name, ".env.")
					fmt.Printf("  %s -> .env + %s\n", suffix, name)
				}
			}
			return nil
		},
	}

	profileCmd.AddCommand(resolveCmd, listCmd)
	rootCmd.AddCommand(profileCmd)
}

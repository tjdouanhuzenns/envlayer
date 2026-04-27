package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envlayer/internal/envfile"
)

func init() {
	var files []string
	var scopes []string
	var format string

	cmd := &cobra.Command{
		Use:   "scope",
		Short: "Merge named scopes from env files in declared order",
		Long: `Load one or more env files as named scopes and resolve them
in order, with later scopes overriding earlier ones.

Example:
  envlayer scope --file base=.env.base --file prod=.env.prod --scopes base,prod`,
		RunE: func(cmd *cobra.Command, args []string) error {
			reg := envfile.NewScopeRegistry()

			for _, entry := range files {
				parts := strings.SplitN(entry, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid --file entry %q: expected name=path", entry)
				}
				name, path := parts[0], parts[1]
				env, err := envfile.ParseFile(path)
				if err != nil {
					return fmt.Errorf("failed to parse %q: %w", path, err)
				}
				if err := reg.Register(name, env); err != nil {
					return fmt.Errorf("failed to register scope %q: %w", name, err)
				}
			}

			if len(scopes) == 0 {
				scopes = reg.List()
			}

			resolved, err := reg.Resolve(scopes...)
			if err != nil {
				return err
			}

			fmt, err := envfile.ParseFormat(format)
			if err != nil {
				return err
			}

			out, err := envfile.Export(resolved, fmt)
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(os.Stdout, out)
			return err
		},
	}

	cmd.Flags().StringArrayVar(&files, "file", nil, "name=path pairs for each scope (repeatable)")
	cmd.Flags().StringSliceVar(&scopes, "scopes", nil, "ordered scope names to resolve (default: all in registration order)")
	cmd.Flags().StringVar(&format, "format", "dotenv", "output format: dotenv, export, json")
	_ = cmd.MarkFlagRequired("file")

	rootCmd.AddCommand(cmd)
}

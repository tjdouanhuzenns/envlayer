package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"envlayer/internal/envfile"
)

func init() {
	var interval int
	var format string

	cmd := &cobra.Command{
		Use:   "watch [files...]",
		Short: "Watch env files for changes and print diffs",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			d := time.Duration(interval) * time.Millisecond
			w := envfile.NewWatcher(d, args...)

			// seed initial state without emitting events
			prev := make(map[string]envfile.EnvMap)
			for _, p := range args {
				env, err := envfile.ParseFile(p)
				if err == nil {
					prev[p] = env
				}
			}

			w.Start()
			defer w.Stop()

			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

			fmt.Fprintf(cmd.OutOrStdout(), "Watching %v (interval %dms)...\n", args, interval)

			for {
				select {
				case ev := <-w.Events:
					fmt.Fprintf(cmd.OutOrStdout(), "\n[changed] %s\n", ev.Path)
					old := prev[ev.Path]
					results := envfile.Diff(old, ev.Env)
					for _, r := range results {
						switch r.Status {
						case "added":
							fmt.Fprintf(cmd.OutOrStdout(), "  + %s=%s\n", r.Key, r.NewValue)
						case "removed":
							fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", r.Key)
						case "changed":
							if format == "verbose" {
								fmt.Fprintf(cmd.OutOrStdout(), "  ~ %s: %q -> %q\n", r.Key, r.OldValue, r.NewValue)
							} else {
								fmt.Fprintf(cmd.OutOrStdout(), "  ~ %s=%s\n", r.Key, r.NewValue)
							}
						}
					}
					prev[ev.Path] = ev.Env
				case err := <-w.Errors:
					log.Printf("watch error: %v", err)
				case <-sig:
					fmt.Fprintln(cmd.OutOrStdout(), "\nStopped.")
					return nil
				}
			}
		},
	}

	cmd.Flags().IntVarP(&interval, "interval", "i", 500, "Poll interval in milliseconds")
	cmd.Flags().StringVarP(&format, "format", "f", "short", "Output format: short|verbose")
	rootCmd.AddCommand(cmd)
}

package main

import (
	"fmt"
	"os"

	"github.com/your-org/envlayer/internal/envfile"
	"github.com/spf13/cobra"
)

func init() {
	var leftName string
	var rightName string

	cmd := &cobra.Command{
		Use:   "compare <left-file> <right-file>",
		Short: "Compare two env files and show differences",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			leftPath := args[0]
			rightPath := args[1]

			if leftName == "" {
				leftName = leftPath
			}
			if rightName == "" {
				rightName = rightPath
			}

			leftMap, err := envfile.ParseFile(leftPath)
			if err != nil {
				return fmt.Errorf("reading %s: %w", leftPath, err)
			}

			rightMap, err := envfile.ParseFile(rightPath)
			if err != nil {
				return fmt.Errorf("reading %s: %w", rightPath, err)
			}

			result := envfile.Compare(leftName, leftMap, rightName, rightMap)

			fmt.Print(result.Summary())

			if len(result.Differ) > 0 || len(result.OnlyLeft) > 0 || len(result.OnlyRight) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&leftName, "left-name", "", "Label for the left file (default: file path)")
	cmd.Flags().StringVar(&rightName, "right-name", "", "Label for the right file (default: file path)")

	rootCmd.AddCommand(cmd)
}

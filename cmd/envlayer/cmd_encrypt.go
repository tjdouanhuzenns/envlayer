package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envlayer/internal/envfile"
)

func init() {
	var passphrase string
	var keys []string
	var decrypt bool
	var output string

	cmd := &cobra.Command{
		Use:   "encrypt <file>",
		Short: "Encrypt or decrypt values in an env file",
		Long: `Encrypt sensitive values in an env file using AES-256-GCM.
Encrypted values are prefixed with "enc:" and base64-encoded.
Use --decrypt to reverse the operation.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if passphrase == "" {
				passphrase = os.Getenv("ENVLAYER_PASSPHRASE")
			}
			if passphrase == "" {
				return fmt.Errorf("passphrase required: use --passphrase or set ENVLAYER_PASSPHRASE")
			}

			env, err := envfile.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("parse %q: %w", args[0], err)
			}

			var result envfile.EnvMap
			if decrypt {
				result, err = envfile.Decrypt(env, passphrase)
				if err != nil {
					return fmt.Errorf("decrypt: %w", err)
				}
			} else {
				var selectedKeys []string
				for _, k := range keys {
					selectedKeys = append(selectedKeys, strings.TrimSpace(k))
				}
				result, err = envfile.Encrypt(env, envfile.EncryptOptions{
					Passphrase: passphrase,
					Keys:       selectedKeys,
				})
				if err != nil {
					return fmt.Errorf("encrypt: %w", err)
				}
			}

			dest := args[0]
			if output != "" {
				dest = output
			}
			if err := envfile.WriteFile(dest, result); err != nil {
				return fmt.Errorf("write %q: %w", dest, err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "wrote %s\n", dest)
			return nil
		},
	}

	cmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase for AES-256-GCM encryption")
	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "Comma-separated keys to encrypt (default: all)")
	cmd.Flags().BoolVarP(&decrypt, "decrypt", "d", false, "Decrypt instead of encrypt")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file (default: overwrite input)")

	rootCmd.AddCommand(cmd)
}

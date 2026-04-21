package envfile

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// ExportFormat defines the output format for environment variable export.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv" // KEY=VALUE
	FormatExport ExportFormat = "export" // export KEY=VALUE
	FormatJSON   ExportFormat = "json"   // {"KEY": "VALUE"}
)

// Export writes an EnvMap to w in the given format.
func Export(w io.Writer, env EnvMap, format ExportFormat) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	switch format {
	case FormatDotenv:
		for _, k := range keys {
			fmt.Fprintf(w, "%s=%s\n", k, quoteIfNeeded(env[k]))
		}
	case FormatExport:
		for _, k := range keys {
			fmt.Fprintf(w, "export %s=%s\n", k, quoteIfNeeded(env[k]))
		}
	case FormatJSON:
		fmt.Fprintln(w, "{")
		for i, k := range keys {
			comma := ","
			if i == len(keys)-1 {
				comma = ""
			}
			fmt.Fprintf(w, "  %q: %q%s\n", k, env[k], comma)
		}
		fmt.Fprintln(w, "}")
	default:
		return fmt.Errorf("unknown export format: %q", format)
	}
	return nil
}

// ExportToFile writes an EnvMap to a file at path in the given format.
func ExportToFile(path string, env EnvMap, format ExportFormat) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating export file: %w", err)
	}
	defer f.Close()
	if err := Export(f, env, format); err != nil {
		return err
	}
	return f.Close()
}

// ParseFormat converts a string to an ExportFormat, returning an error if the
// value is not a recognised format name.
func ParseFormat(s string) (ExportFormat, error) {
	switch ExportFormat(strings.ToLower(s)) {
	case FormatDotenv:
		return FormatDotenv, nil
	case FormatExport:
		return FormatExport, nil
	case FormatJSON:
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unknown export format: %q (valid: dotenv, export, json)", s)
	}
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t\n#") {
		return fmt.Sprintf("%q", v)
	}
	return v
}

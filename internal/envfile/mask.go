package envfile

import "strings"

// MaskOptions controls how sensitive values are masked.
type MaskOptions struct {
	// Keys is a list of exact key names to mask.
	Keys []string
	// Patterns is a list of substrings; any key containing one will be masked.
	Patterns []string
	// MaskChar is the character used for masking. Defaults to "*".
	MaskChar string
	// VisibleChars is how many trailing characters to leave visible. Defaults to 0.
	VisibleChars int
}

// defaultPatterns are common sensitive key substrings.
var defaultPatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "APIKEY", "PRIVATE", "CREDENTIAL",
}

// Mask returns a copy of env with sensitive values replaced by a mask string.
// If opts is nil, default sensitive patterns are applied.
func Mask(env EnvMap, opts *MaskOptions) EnvMap {
	if opts == nil {
		opts = &MaskOptions{Patterns: defaultPatterns}
	}
	maskChar := opts.MaskChar
	if maskChar == "" {
		maskChar = "*"
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	result := make(EnvMap, len(env))
	for k, v := range env {
		if shouldMask(k, keySet, opts.Patterns) {
			result[k] = maskValue(v, maskChar, opts.VisibleChars)
		} else {
			result[k] = v
		}
	}
	return result
}

func shouldMask(key string, keySet map[string]struct{}, patterns []string) bool {
	if _, ok := keySet[key]; ok {
		return true
	}
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

func maskValue(value, maskChar string, visibleChars int) string {
	if len(value) == 0 {
		return ""
	}
	if visibleChars <= 0 || visibleChars >= len(value) {
		return strings.Repeat(maskChar, 6)
	}
	visible := value[len(value)-visibleChars:]
	return strings.Repeat(maskChar, 6) + visible
}

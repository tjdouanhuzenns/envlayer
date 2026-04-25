package envfile

import (
	"fmt"
	"strings"
)

// LintRule represents a single linting check.
type LintRule string

const (
	LintRuleDuplicateKey    LintRule = "duplicate_key"
	LintRuleEmptyValue      LintRule = "empty_value"
	LintRuleAllCaps         LintRule = "not_all_caps"
	LintRuleLeadingUnderscore LintRule = "leading_underscore"
)

// LintIssue describes a single lint finding.
type LintIssue struct {
	Key     string
	Rule    LintRule
	Message string
}

func (i LintIssue) String() string {
	return fmt.Sprintf("[%s] %s: %s", i.Rule, i.Key, i.Message)
}

// LintOptions controls which rules are active.
type LintOptions struct {
	WarnEmptyValues      bool
	WarnNotAllCaps       bool
	WarnLeadingUnderscore bool
}

// DefaultLintOptions returns sensible defaults.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		WarnEmptyValues:      true,
		WarnNotAllCaps:       true,
		WarnLeadingUnderscore: false,
	}
}

// Lint checks an EnvMap for style and quality issues.
func Lint(env EnvMap, opts LintOptions) []LintIssue {
	var issues []LintIssue

	for key, val := range env {
		if opts.WarnEmptyValues && strings.TrimSpace(val) == "" {
			issues = append(issues, LintIssue{
				Key:     key,
				Rule:    LintRuleEmptyValue,
				Message: "value is empty or whitespace-only",
			})
		}

		if opts.WarnNotAllCaps && key != strings.ToUpper(key) {
			issues = append(issues, LintIssue{
				Key:     key,
				Rule:    LintRuleAllCaps,
				Message: fmt.Sprintf("key %q is not all uppercase", key),
			})
		}

		if opts.WarnLeadingUnderscore && strings.HasPrefix(key, "_") {
			issues = append(issues, LintIssue{
				Key:     key,
				Rule:    LintRuleLeadingUnderscore,
				Message: fmt.Sprintf("key %q starts with an underscore", key),
			})
		}
	}

	return issues
}

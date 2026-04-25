package envfile

import (
	"testing"
)

func TestLint_EmptyValue(t *testing.T) {
	env := EnvMap{"HOST": "localhost", "EMPTY": ""}
	opts := DefaultLintOptions()
	issues := Lint(env, opts)
	if !hasIssueForKey(issues, "EMPTY", LintRuleEmptyValue) {
		t.Error("expected empty value issue for EMPTY")
	}
	if hasIssueForKey(issues, "HOST", LintRuleEmptyValue) {
		t.Error("unexpected empty value issue for HOST")
	}
}

func TestLint_NotAllCaps(t *testing.T) {
	env := EnvMap{"myKey": "value", "GOOD_KEY": "value"}
	opts := DefaultLintOptions()
	issues := Lint(env, opts)
	if !hasIssueForKey(issues, "myKey", LintRuleAllCaps) {
		t.Error("expected not_all_caps issue for myKey")
	}
	if hasIssueForKey(issues, "GOOD_KEY", LintRuleAllCaps) {
		t.Error("unexpected not_all_caps issue for GOOD_KEY")
	}
}

func TestLint_LeadingUnderscore(t *testing.T) {
	env := EnvMap{"_PRIVATE": "secret", "PUBLIC": "yes"}
	opts := DefaultLintOptions()
	opts.WarnLeadingUnderscore = true
	issues := Lint(env, opts)
	if !hasIssueForKey(issues, "_PRIVATE", LintRuleLeadingUnderscore) {
		t.Error("expected leading_underscore issue for _PRIVATE")
	}
	if hasIssueForKey(issues, "PUBLIC", LintRuleLeadingUnderscore) {
		t.Error("unexpected leading_underscore issue for PUBLIC")
	}
}

func TestLint_NoIssues(t *testing.T) {
	env := EnvMap{"APP_HOST": "localhost", "APP_PORT": "8080"}
	opts := DefaultLintOptions()
	issues := Lint(env, opts)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLint_DisabledRules(t *testing.T) {
	env := EnvMap{"lowercase": ""}
	opts := LintOptions{
		WarnEmptyValues: false,
		WarnNotAllCaps:  false,
	}
	issues := Lint(env, opts)
	if len(issues) != 0 {
		t.Errorf("expected no issues with all rules disabled, got %d", len(issues))
	}
}

func TestLintIssue_String(t *testing.T) {
	issue := LintIssue{Key: "FOO", Rule: LintRuleEmptyValue, Message: "value is empty"}
	s := issue.String()
	if s == "" {
		t.Error("expected non-empty string from LintIssue.String()")
	}
}

// helper
func hasIssueForKey(issues []LintIssue, key string, rule LintRule) bool {
	for _, i := range issues {
		if i.Key == key && i.Rule == rule {
			return true
		}
	}
	return false
}

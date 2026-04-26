package envfile

import (
	"strings"
	"testing"
)

func TestCompare_OnlyLeft(t *testing.T) {
	left := EnvMap{"A": "1", "B": "2"}
	right := EnvMap{"A": "1"}
	r := Compare("dev", left, "prod", right)
	if len(r.OnlyLeft) != 1 || r.OnlyLeft[0] != "B" {
		t.Errorf("expected OnlyLeft=[B], got %v", r.OnlyLeft)
	}
	if len(r.OnlyRight) != 0 {
		t.Errorf("expected empty OnlyRight, got %v", r.OnlyRight)
	}
}

func TestCompare_OnlyRight(t *testing.T) {
	left := EnvMap{"A": "1"}
	right := EnvMap{"A": "1", "C": "3"}
	r := Compare("dev", left, "prod", right)
	if len(r.OnlyRight) != 1 || r.OnlyRight[0] != "C" {
		t.Errorf("expected OnlyRight=[C], got %v", r.OnlyRight)
	}
}

func TestCompare_Differ(t *testing.T) {
	left := EnvMap{"HOST": "localhost", "PORT": "8080"}
	right := EnvMap{"HOST": "prod.example.com", "PORT": "8080"}
	r := Compare("dev", left, "prod", right)
	if len(r.Differ) != 1 {
		t.Fatalf("expected 1 difference, got %d", len(r.Differ))
	}
	pair, ok := r.Differ["HOST"]
	if !ok {
		t.Fatal("expected HOST in Differ")
	}
	if pair[0] != "localhost" || pair[1] != "prod.example.com" {
		t.Errorf("unexpected pair: %v", pair)
	}
}

func TestCompare_Same(t *testing.T) {
	left := EnvMap{"A": "1", "B": "2"}
	right := EnvMap{"A": "1", "B": "2"}
	r := Compare("dev", left, "prod", right)
	if len(r.Same) != 2 {
		t.Errorf("expected 2 same keys, got %d", len(r.Same))
	}
	if len(r.Differ) != 0 || len(r.OnlyLeft) != 0 || len(r.OnlyRight) != 0 {
		t.Error("expected no differences")
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	r := Compare("dev", EnvMap{}, "prod", EnvMap{})
	if len(r.OnlyLeft)+len(r.OnlyRight)+len(r.Differ)+len(r.Same) != 0 {
		t.Error("expected empty result for empty maps")
	}
}

func TestCompare_Summary(t *testing.T) {
	left := EnvMap{"HOST": "localhost", "ONLY_DEV": "yes"}
	right := EnvMap{"HOST": "prod.host", "ONLY_PROD": "yes"}
	r := Compare("dev", left, "prod", right)
	summary := r.Summary()
	if !strings.Contains(summary, "dev") || !strings.Contains(summary, "prod") {
		t.Error("summary should contain env names")
	}
	if !strings.Contains(summary, "ONLY_DEV") {
		t.Error("summary should mention ONLY_DEV")
	}
	if !strings.Contains(summary, "HOST") {
		t.Error("summary should mention HOST")
	}
}

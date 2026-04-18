package envfile

import "testing"

func makeEnv(pairs ...string) *EnvMap {
	e := NewEnvMap()
	for i := 0; i+1 < len(pairs); i += 2 {
		e.Set(pairs[i], pairs[i+1])
	}
	return e
}

func TestMerge_Override(t *testing.T) {
	base := makeEnv("APP_ENV", "dev", "PORT", "3000")
	layer := makeEnv("APP_ENV", "prod", "SECRET", "xyz")

	result, err := Merge(base, []*EnvMap{layer}, StrategyOverride)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Values["APP_ENV"] != "prod" {
		t.Errorf("expected prod, got %s", result.Values["APP_ENV"])
	}
	if result.Values["PORT"] != "3000" {
		t.Errorf("expected 3000, got %s", result.Values["PORT"])
	}
	if result.Values["SECRET"] != "xyz" {
		t.Errorf("expected xyz, got %s", result.Values["SECRET"])
	}
}

func TestMerge_KeepBase(t *testing.T) {
	base := makeEnv("APP_ENV", "dev", "PORT", "3000")
	layer := makeEnv("APP_ENV", "prod", "SECRET", "xyz")

	result, err := Merge(base, []*EnvMap{layer}, StrategyKeepBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Values["APP_ENV"] != "dev" {
		t.Errorf("expected dev (base preserved), got %s", result.Values["APP_ENV"])
	}
	if result.Values["SECRET"] != "xyz" {
		t.Errorf("expected xyz (new key added), got %s", result.Values["SECRET"])
	}
}

func TestMerge_NilBase(t *testing.T) {
	_, err := Merge(nil, nil, StrategyOverride)
	if err == nil {
		t.Fatal("expected error for nil base")
	}
}

func TestMerge_NilLayerSkipped(t *testing.T) {
	base := makeEnv("KEY", "val")
	result, err := Merge(base, []*EnvMap{nil}, StrategyOverride)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Values["KEY"] != "val" {
		t.Errorf("expected val, got %s", result.Values["KEY"])
	}
}

func TestMergeFiles(t *testing.T) {
	base := writeTempEnv(t, "APP_ENV=dev\nPORT=3000\n")
	layer := writeTempEnv(t, "APP_ENV=staging\nDEBUG=true\n")

	result, err := MergeFiles([]string{base, layer}, StrategyOverride)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Values["APP_ENV"] != "staging" {
		t.Errorf("expected staging, got %s", result.Values["APP_ENV"])
	}
	if result.Values["DEBUG"] != "true" {
		t.Errorf("expected true, got %s", result.Values["DEBUG"])
	}
}

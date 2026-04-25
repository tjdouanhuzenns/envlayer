package envfile

import (
	"strings"
	"testing"
)

func TestProfile_ResolveAndExport(t *testing.T) {
	dir := t.TempDir()
	base := writeProfileEnv(t, dir, ".env", "APP=myapp\nLOG_LEVEL=info\nSECRET_KEY=base-secret\n")
	prod := writeProfileEnv(t, dir, ".env.prod", "LOG_LEVEL=warn\nSECRET_KEY=prod-secret\nREGION=eu-west\n")

	reg := NewProfileRegistry()
	_ = reg.Register(&Profile{Name: "production", Layers: []string{base, prod}})

	env, err := reg.Resolve("production")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}

	out, err := Export(env, FormatDotenv)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if !strings.Contains(out, "APP=myapp") {
		t.Errorf("missing APP=myapp in output")
	}
	if !strings.Contains(out, "LOG_LEVEL=warn") {
		t.Errorf("expected LOG_LEVEL=warn override")
	}
	if !strings.Contains(out, "SECRET_KEY=prod-secret") {
		t.Errorf("expected SECRET_KEY=prod-secret override")
	}
	if !strings.Contains(out, "REGION=eu-west") {
		t.Errorf("expected REGION=eu-west from prod layer")
	}
}

func TestProfile_ResolveAndMask(t *testing.T) {
	dir := t.TempDir()
	base := writeProfileEnv(t, dir, ".env", "APP=myapp\nSECRET_KEY=supersecret\nAPI_TOKEN=token123\n")

	reg := NewProfileRegistry()
	_ = reg.Register(&Profile{Name: "dev", Layers: []string{base}})

	env, err := reg.Resolve("dev")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}

	masked := Mask(env, MaskOptions{})

	if masked["APP"] != "myapp" {
		t.Errorf("APP should not be masked, got %s", masked["APP"])
	}
	if masked["SECRET_KEY"] == "supersecret" {
		t.Errorf("SECRET_KEY should be masked")
	}
	if masked["API_TOKEN"] == "token123" {
		t.Errorf("API_TOKEN should be masked")
	}
}

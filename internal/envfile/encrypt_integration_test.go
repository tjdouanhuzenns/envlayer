package envfile

import (
	"os"
	"testing"
)

// TestEncrypt_WriteAndReadBack encrypts an env map, writes it to a file,
// reads it back, and decrypts — verifying full round-trip through the FS.
func TestEncrypt_WriteAndReadBack(t *testing.T) {
	original := EnvMap{
		"DB_PASSWORD": "s3cr3t!",
		"API_KEY":     "abc-123-xyz",
		"APP_ENV":     "production",
	}

	encrypted, err := Encrypt(original, EncryptOptions{
		Passphrase: "integration-pass",
		Keys:       []string{"DB_PASSWORD", "API_KEY"},
	})
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	tmp, err := os.CreateTemp("", "envlayer-enc-*.env")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	tmp.Close()

	if err := WriteFile(tmp.Name(), encrypted); err != nil {
		t.Fatalf("write: %v", err)
	}

	loaded, err := ParseFile(tmp.Name())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	decrypted, err := Decrypt(loaded, "integration-pass")
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}

	for k, want := range original {
		if got := decrypted[k]; got != want {
			t.Errorf("key %q: want %q, got %q", k, want, got)
		}
	}

	// APP_ENV was not in the Keys list — should remain plain
	if encrypted["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should remain plain, got %q", encrypted["APP_ENV"])
	}
}

// TestEncrypt_ThenMask ensures encrypted values can be masked before display.
func TestEncrypt_ThenMask(t *testing.T) {
	env := EnvMap{"SECRET_KEY": "plaintext", "NAME": "app"}
	enc, err := Encrypt(env, EncryptOptions{Passphrase: "pass", Keys: []string{"SECRET_KEY"}})
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	masked := Mask(enc, MaskOptions{})
	if masked["SECRET_KEY"] == "plaintext" {
		t.Error("SECRET_KEY should be masked")
	}
}

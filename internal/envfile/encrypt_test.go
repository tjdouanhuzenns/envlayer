package envfile

import (
	"strings"
	"testing"
)

func TestEncrypt_AllKeys(t *testing.T) {
	env := EnvMap{"SECRET": "topsecret", "NAME": "app"}
	enc, err := Encrypt(env, EncryptOptions{Passphrase: "pass123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range enc {
		if !strings.HasPrefix(v, "enc:") {
			t.Errorf("key %q not encrypted: %q", k, v)
		}
	}
}

func TestEncrypt_SelectedKeys(t *testing.T) {
	env := EnvMap{"SECRET": "topsecret", "NAME": "app"}
	enc, err := Encrypt(env, EncryptOptions{Passphrase: "pass", Keys: []string{"SECRET"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(enc["SECRET"], "enc:") {
		t.Errorf("SECRET should be encrypted")
	}
	if enc["NAME"] != "app" {
		t.Errorf("NAME should be unchanged, got %q", enc["NAME"])
	}
}

func TestDecrypt_RoundTrip(t *testing.T) {
	env := EnvMap{"DB_PASS": "hunter2", "HOST": "localhost"}
	enc, err := Encrypt(env, EncryptOptions{Passphrase: "secret"})
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}
	dec, err := Decrypt(enc, "secret")
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}
	for k, want := range env {
		if got := dec[k]; got != want {
			t.Errorf("key %q: want %q, got %q", k, want, got)
		}
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	env := EnvMap{"X": "value"}
	enc, _ := Encrypt(env, EncryptOptions{Passphrase: "correct"})
	_, err := Decrypt(enc, "wrong")
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}

func TestEncrypt_EmptyPassphrase(t *testing.T) {
	env := EnvMap{"K": "v"}
	_, err := Encrypt(env, EncryptOptions{Passphrase: ""})
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestDecrypt_EmptyPassphrase(t *testing.T) {
	env := EnvMap{"K": "v"}
	_, err := Decrypt(env, "")
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestEncrypt_AlreadyEncryptedSkipped(t *testing.T) {
	env := EnvMap{"K": "enc:alreadyencoded"}
	enc, err := Encrypt(env, EncryptOptions{Passphrase: "pass"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enc["K"] != "enc:alreadyencoded" {
		t.Errorf("already-encrypted value should not be re-encrypted")
	}
}

func TestDecrypt_PlainValuesPassThrough(t *testing.T) {
	env := EnvMap{"PLAIN": "hello"}
	dec, err := Decrypt(env, "anypass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec["PLAIN"] != "hello" {
		t.Errorf("plain value should pass through unchanged")
	}
}

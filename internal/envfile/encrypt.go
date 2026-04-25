package envfile

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// EncryptOptions configures encryption behavior.
type EncryptOptions struct {
	// Keys lists the env var keys to encrypt. If empty, all keys are encrypted.
	Keys []string
	// Passphrase is used to derive the AES-256 key.
	Passphrase string
}

// Encrypt encrypts values in the env map using AES-256-GCM.
// Encrypted values are base64-encoded and prefixed with "enc:".
func Encrypt(env EnvMap, opts EncryptOptions) (EnvMap, error) {
	if opts.Passphrase == "" {
		return nil, errors.New("encrypt: passphrase must not be empty")
	}
	key := deriveKey(opts.Passphrase)
	targets := resolveTargetKeys(env, opts.Keys)

	result := make(EnvMap, len(env))
	for k, v := range env {
		result[k] = v
	}

	for _, k := range targets {
		v, ok := result[k]
		if !ok {
			continue
		}
		if isEncrypted(v) {
			continue
		}
		enc, err := aesGCMEncrypt(key, []byte(v))
		if err != nil {
			return nil, fmt.Errorf("encrypt: key %q: %w", k, err)
		}
		result[k] = "enc:" + base64.StdEncoding.EncodeToString(enc)
	}
	return result, nil
}

// Decrypt decrypts values in the env map that are prefixed with "enc:".
func Decrypt(env EnvMap, passphrase string) (EnvMap, error) {
	if passphrase == "" {
		return nil, errors.New("decrypt: passphrase must not be empty")
	}
	key := deriveKey(passphrase)
	result := make(EnvMap, len(env))
	for k, v := range env {
		if !isEncrypted(v) {
			result[k] = v
			continue
		}
		ciphertext, err := base64.StdEncoding.DecodeString(v[4:])
		if err != nil {
			return nil, fmt.Errorf("decrypt: key %q: invalid base64: %w", k, err)
		}
		plain, err := aesGCMDecrypt(key, ciphertext)
		if err != nil {
			return nil, fmt.Errorf("decrypt: key %q: %w", k, err)
		}
		result[k] = string(plain)
	}
	return result, nil
}

func isEncrypted(v string) bool {
	return len(v) > 4 && v[:4] == "enc:"
}

func deriveKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func resolveTargetKeys(env EnvMap, keys []string) []string {
	if len(keys) > 0 {
		return keys
	}
	out := make([]string, 0, len(env))
	for k := range env {
		out = append(out, k)
	}
	return out
}

func aesGCMEncrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func aesGCMDecrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ns := gcm.NonceSize()
	if len(ciphertext) < ns {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ct := ciphertext[:ns], ciphertext[ns:]
	return gcm.Open(nil, nonce, ct, nil)
}

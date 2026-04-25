// Package envfile provides encrypt/decrypt support for env maps.
//
// # Encrypt
//
// Encrypt encrypts selected (or all) values in an EnvMap using AES-256-GCM.
// The passphrase is hashed with SHA-256 to produce a 32-byte key.
// Encrypted values are stored as "enc:<base64>" so they remain valid in
// standard .env files and can be committed safely.
//
// Example:
//
//	env := envfile.EnvMap{"DB_PASSWORD": "s3cr3t", "APP_NAME": "myapp"}
//	encrypted, err := envfile.Encrypt(env, envfile.EncryptOptions{
//		Passphrase: "my-passphrase",
//		Keys:       []string{"DB_PASSWORD"},
//	})
//
// # Decrypt
//
// Decrypt reverses the process. Only values prefixed with "enc:" are
// decrypted; plain values are passed through unchanged.
//
//	decrypted, err := envfile.Decrypt(encrypted, "my-passphrase")
package envfile

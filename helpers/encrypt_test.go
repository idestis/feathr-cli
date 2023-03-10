package helpers

import (
	"testing"
)

func TestEncryptAndDecryptString(t *testing.T) {
	plaintext := "Hello, world!"

	// Encrypt plaintext
	ciphertext, err := EncryptString(plaintext)
	if err != nil {
		t.Errorf("Error encrypting string: %v", err)
	}

	// Decrypt ciphertext
	decrypted, err := DecryptString(ciphertext)
	if err != nil {
		t.Errorf("Error decrypting string: %v", err)
	}

	// Check that decrypted plaintext matches original plaintext
	if decrypted != plaintext {
		t.Errorf("Decrypted string %q does not match original plaintext %q", decrypted, plaintext)
	}
}

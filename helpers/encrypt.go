package helpers

import (
	"encoding/base64"
)

// EncryptString does a simple Base64 encryption on a string
// One day I'll write a better encryption algorithm using AES
func EncryptString(text string) (string, error) {
	plaintext := []byte(text)
	key := []byte("d6#mX9^Lfeathr5Y@U$%F8Cgcli#&K7p+Z*2hE!")

	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i++ {
		ciphertext[i] = plaintext[i] ^ key[i%len(key)]
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// EncryptString does a simple Base64 decryption on a string
// One day I'll write a better encryption algorithm using AES
func DecryptString(ciphertext string) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	key := []byte("d6#mX9^Lfeathr5Y@U$%F8Cgcli#&K7p+Z*2hE!") // I know, this is not a good key

	plaintext := make([]byte, len(ciphertextBytes))
	for i := 0; i < len(ciphertextBytes); i++ {
		plaintext[i] = ciphertextBytes[i] ^ key[i%len(key)]
	}

	return string(plaintext), nil
}

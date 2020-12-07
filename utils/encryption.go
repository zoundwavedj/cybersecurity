package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"os"
)

var (
	encodedKey = os.Getenv("ENCRYPT_KEY")
)

func initCipher() (cipher.AEAD, error) {
	key, err := hex.DecodeString(encodedKey)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}

// Encrypt function
func Encrypt(plaintext string) (string, error) {
	cipher, err := initCipher()
	if err != nil {
		return "", err
	}

	nonce := make([]byte, cipher.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := cipher.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.RawStdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt function
func Decrypt(ciphertext string) (string, error) {
	cipher, err := initCipher()
	if err != nil {
		return "", err
	}

	decodedCipherText, err := base64.RawStdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	nonceSize := cipher.NonceSize()
	nonce, encryptedstring := decodedCipherText[:nonceSize], decodedCipherText[nonceSize:]

	plaintext, err := cipher.Open(nil, nonce, encryptedstring, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

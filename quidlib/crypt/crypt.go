package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/teal-finance/quid/quidlib/conf"
)

const (
	// nonceSize= 12 // AES-128 nonce is 12 bytes
	gcmTagSize = 16 // AES-GCM tag is 16 bytes
)

// AesGcmEncryptHex : encrypt content.
func AesGcmEncryptHex(plaintext string) (string, error) {
	b, err := AesGcmEncryptBin([]byte(plaintext))
	return hex.EncodeToString(b), err
}

// AesGcmEncrypt : encrypt content.
func AesGcmEncryptBin(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(conf.EncodingKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	returnedSize := gcm.NonceSize() + len(plaintext) + gcmTagSize

	iv := make([]byte, gcm.NonceSize(), returnedSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("random iv generation: %w", err)
	}

	ciphertext := gcm.Seal(nil, iv, plaintext, nil)
	return append(iv, ciphertext...), nil
}

// AesGcmDecryptHex : decrypt content.
func AesGcmDecryptHex(encryptedString string) (string, error) {
	bytes, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	plaintext, err := AesGcmDecryptBin(bytes)
	return string(plaintext), err
}

func AesGcmDecryptBin(bytes []byte) ([]byte, error) {
	block, err := aes.NewCipher(conf.EncodingKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	iv, ciphertext := bytes[:gcm.NonceSize()], bytes[gcm.NonceSize():]

	// we are not subject to confused deputy attack => additionalData can be empty
	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

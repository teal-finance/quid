package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// EncodingKey is used to encode each JWT secret key in the DB.
var EncodingKey []byte

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
	block, err := aes.NewCipher(EncodingKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// all will contain the nonce (first 12 bytes) + the cypher text + the GCM tag
	all := make([]byte, gcm.NonceSize(), gcm.NonceSize()+len(plaintext)+gcmTagSize)

	// write a random nonce
	if _, err := rand.Read(all); err != nil {
		return nil, fmt.Errorf("random iv generation: %w", err)
	}

	// write the cypher text after the nonce and appends the GCM tag
	all = gcm.Seal(all, all, plaintext, nil)
	return all, nil
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
	block, err := aes.NewCipher(EncodingKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	iv := encryptedBytes[:gcm.NonceSize()]
	ciphertext := encryptedBytes[gcm.NonceSize():]
	dst := ciphertext[:0]

	// we are not subject to confused deputy attack => additionalData can be empty
	plaintext, err := gcm.Open(dst, iv, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

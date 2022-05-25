package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/teal-finance/quid/quidlib/conf"

	rand "github.com/zhangyunhao116/fastrand"
)

const (
	aesGcmNonceSize = 12
)

// AesGcmEncrypt : encrypt content.
func AesGcmEncrypt(plaintext string, additionalData []byte) (string, error) {
	key, err := hex.DecodeString(conf.EncodingKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aesGcmNonceSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("random iv generation : %w", err)
	}

	ciphertext := aesgcm.Seal(nil, iv, []byte(plaintext), additionalData)

	return hex.EncodeToString(append(iv, ciphertext...)), nil
}

// AesGcmDecrypt : decrypt content.
func AesGcmDecrypt(encryptedString string, additionalData []byte) (string, error) {
	key, err := hex.DecodeString(conf.EncodingKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	iv, ciphertext := enc[:aesGcmNonceSize], enc[aesGcmNonceSize:]

	plaintext, err := aesgcm.Open(nil, iv, ciphertext, additionalData)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

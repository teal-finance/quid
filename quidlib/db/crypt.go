package db

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"

	"github.com/synw/quid/quidlib/conf"
)

func aesGcmEncrypt(plaintext string, additionalData []byte) (string, error) {
	key, _ := hex.DecodeString(conf.EncodingKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	iv := make([]byte, aesgcm.NonceSize())
	ciphertext := aesgcm.Seal(iv, iv, []byte(plaintext), additionalData)
	return hex.EncodeToString(ciphertext), nil
}

func aesGcmDecrypt(encryptedString string, additionalData []byte) (string, error) {
	key, _ := hex.DecodeString(conf.EncodingKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesgcm.NonceSize()
	enc, _ := hex.DecodeString(encryptedString)
	iv, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesgcm.Open(nil, iv, ciphertext, additionalData)
	if err != nil {
		return "", err
	}
	s := string(plaintext[:])
	return s, nil
}

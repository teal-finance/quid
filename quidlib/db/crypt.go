package db

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

func encrypt(plaintext, key string) (string, error) {
	k := []byte(key)
	c, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	ciphertext := gcmEncrypt(key, plaintext, nonce, nil)
	return ciphertext, nil
}

func decrypt(plaintext, key string) (string, error) {
	k := []byte(key)
	c, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	ciphertext := gcmDecrypt(key, plaintext, nonce, nil)
	return ciphertext, nil
}

func gcmEncrypt(key string, plaintext string, iv []byte, additionalData []byte) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	ciphertext := aesgcm.Seal(nil, iv, []byte(plaintext), additionalData)
	return hex.EncodeToString(ciphertext)
}

func gcmDecrypt(key string, ct string, iv []byte, additionalData []byte) string {
	ciphertext, _ := hex.DecodeString(ct)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	plaintext, err := aesgcm.Open(nil, iv, ciphertext, additionalData)
	if err != nil {
		panic(err.Error())
	}
	s := string(plaintext[:])
	return s
}

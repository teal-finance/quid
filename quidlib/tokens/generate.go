package tokens

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenUserToken : generate a token for a user
func GenUserToken(name, key string, groups []string, timeout string) (string, error) {
	d, err := time.ParseDuration(timeout)
	if err != nil {
		return "", err
	}
	to := time.Now().Add(d)
	claims := standardUserClaims(name, groups, to)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return token, nil
}

// GenAdminToken : generate a token for a quid admin frontend user
func GenAdminToken(name, key string) (string, error) {
	timeout := time.Now().Add(time.Hour * 24)
	claims := standardUserClaims(name, []string{"quid_admin"}, timeout)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return token, nil
}

// GenKey : generate a random key
func GenKey() string {
	b, err := generateRandomBytes(32)
	if err != nil {
		log.Fatal(err)
	}
	h := hmac.New(sha256.New, []byte(b))
	token := hex.EncodeToString(h.Sum(nil))
	return token

}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

/*
func genKeyForNamespace(name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ns_name": name,
	})
	tokenString, err := token.SignedString(conf.EncodingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}*/

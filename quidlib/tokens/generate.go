package tokens

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/karrick/tparse"
)

// GenUserToken : generate a token for a user
func GenUserToken(name, key string, groups []string, timeout, maxTimeout string) (bool, string, error) {
	isAuthorized, err := isTimeoutAuthorized(timeout, maxTimeout)
	if err != nil {
		emo.ParamError(err)
		return false, "", err
	}
	if !isAuthorized {
		return false, "", nil
	}
	to, err := tparse.ParseNow(time.RFC3339, "now+"+timeout)
	if err != nil {
		emo.TimeError(err)
		return false, "", err
	}
	claims := standardUserClaims(name, groups, to)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(key))
	if err != nil {
		emo.EncryptError(err)
		return false, "", err
	}
	return true, token, nil
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

// GenKey : generate a random hmac key
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

func isTimeoutAuthorized(timeout, maxTimeout string) (bool, error) {
	requested, err := tparse.ParseNow(time.RFC3339, "now+"+timeout)
	if err != nil {
		return false, err
	}
	max, err := tparse.ParseNow(time.RFC3339, "now+"+maxTimeout)
	if err != nil {
		return false, err
	}
	return requested.Before(max), err
}

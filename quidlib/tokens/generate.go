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

// GenRefreshToken : generate a refresh token for a user in a namespace
func GenRefreshToken(namespaceName, namespaceRefreshKey, maxRefreshokenTTL, username string, timeout string) (bool, string, error) {
	isAuthorized, err := isTimeoutAuthorized(timeout, maxRefreshokenTTL)
	if err != nil {
		emo.ParamError(err)
		return false, "", err
	}
	if !isAuthorized {
		emo.ParamError("Unauthorized timeout", timeout)
		return false, "", nil
	}
	to, err := tparse.ParseNow(time.RFC3339, "now+"+timeout)
	claims := standardRefreshClaims(namespaceName, username, to)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(namespaceRefreshKey))
	if err != nil {
		emo.EncryptError(err)
		return false, "", err
	}
	return true, token, nil
}

// GenAccessToken : generate an access token for a user in a namespace
func GenAccessToken(namespaceName, namespaceKey, maxTokenTTL, name string, groups []string, timeout string) (bool, string, error) {
	isAuthorized, err := isTimeoutAuthorized(timeout, maxTokenTTL)
	if err != nil {
		emo.ParamError(err)
		return false, "", err
	}
	if !isAuthorized {
		emo.ParamError("Unauthorized timeout", timeout)
		return false, "", nil
	}
	to, err := tparse.ParseNow(time.RFC3339, "now+"+timeout)
	if err != nil {
		emo.TimeError(err)
		return false, "", err
	}
	claims := standardAccessClaims(namespaceName, name, groups, to)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(namespaceKey))
	if err != nil {
		emo.EncryptError(err)
		return false, "", err
	}
	return true, token, nil
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
	max, err := tparse.ParseNow(time.RFC3339, "now+1s+"+maxTimeout)
	if err != nil {
		return false, err
	}
	return requested.UTC().Before(max), err
}

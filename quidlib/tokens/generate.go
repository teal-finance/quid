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
	"github.com/synw/quid/quidlib/models"
)

// GenRefreshToken : generate a refresh token for a user in a namespace
func GenRefreshToken(namespace models.Namespace, username string, timeout string) (bool, string, error) {
	isAuthorized, err := isTimeoutAuthorized(timeout, namespace.MaxRefreshTokenTTL)
	if err != nil {
		emo.ParamError(err)
		return false, "", err
	}
	if !isAuthorized {
		return false, "", nil
	}
	to, err := tparse.ParseNow(time.RFC3339, "now+"+timeout)
	claims := standardRefreshClaims(namespace.Name, username, to)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(namespace.RefreshKey))
	if err != nil {
		emo.EncryptError(err)
		return false, "", err
	}
	return true, token, nil
}

// GenAccessToken : generate an access token for a user in a namespace
func GenAccessToken(namespace models.Namespace, name string, groups []string, timeout, maxTimeout string) (bool, string, error) {
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
	claims := standardAccessClaims(namespace.Name, name, groups, to)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(namespace.Key))
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
	max, err := tparse.ParseNow(time.RFC3339, "now+"+maxTimeout)
	if err != nil {
		return false, err
	}
	return requested.Before(max), err
}

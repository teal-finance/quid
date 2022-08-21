package tokens

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/karrick/tparse/v2"
)

// GenRefreshToken generates a refresh token for a user in a namespace.
func GenRefreshToken(timeout, maxTTL, namespace, user string, secretKey []byte) (string, error) {
	isAuthorized, err := isTimeoutAuthorized(timeout, maxTTL)
	if err != nil {
		emo.ParamError(err)
		return "", err
	}

	if !isAuthorized {
		emo.ParamError("Unauthorized timeout", timeout)
		return "", nil
	}

	to, err := tparse.ParseNow(time.RFC3339, "now+"+timeout)
	if err != nil {
		emo.ParamError(err)
		return "", err
	}

	claims := newRefreshClaims(namespace, user, to.UTC())
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString(secretKey)
	if err != nil {
		emo.EncryptError(err)
		return "", err
	}

	emo.RefreshToken("Issued a refresh token for user '" + user + "' and namespace " + namespace)
	return token, nil
}

// GenAdminAccessToken generates an admin access token for a user.
func GenAdminAccessToken(namespaceName, timeout, maxTTL, userName string, userId, nsId int64, secretKey []byte, isAdmin, isNsAdmin bool) (string, error) {
	isAuthorized, err := isTimeoutAuthorized(timeout, maxTTL)
	if err != nil {
		emo.ParamError(err)
		return "", err
	}

	if !isAuthorized {
		emo.ParamError("Unauthorized timeout", timeout)
		return "", nil
	}

	to, err := tparse.ParseNow(time.RFC3339, "now+"+timeout)
	if err != nil {
		emo.ParamError(err)
		return "", err
	}

	claims := newAdminAccessClaims(namespaceName, userName, userId, nsId, to.UTC(), isAdmin, isNsAdmin)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString(secretKey)
	if err != nil {
		emo.EncryptError(err)
		return "", err
	}

	emo.AccessToken("Issued an admin access token for user", userName, "and namespace", namespaceName)
	return token, nil
}

// GenAccessToken generates an access token for a user.
func GenAccessToken(timeout, maxTTL, user string, groups, orgs []string, secretKey []byte) (string, error) {
	isAuthorized, err := isTimeoutAuthorized(timeout, maxTTL)
	if err != nil {
		emo.ParamError(err)
		return "", err
	}

	if !isAuthorized {
		emo.ParamError("Unauthorized timeout", timeout)
		return "", nil
	}

	to, err := tparse.ParseNow(time.RFC3339, "now+"+timeout)
	if err != nil {
		emo.ParamError(err)
		return "", err
	}

	claims := newAccessClaims(user, groups, orgs, to.UTC())
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString(secretKey)
	if err != nil {
		emo.EncryptError(err)
		return "", err
	}

	emo.AccessToken("Issued an access token for user", user)

	return token, nil
}

// GenKey generates a random hmac key.
func GenKey() string {
	b := genRandomBytes(32)
	h := hmac.New(sha256.New, b)
	return hex.EncodeToString(h.Sum(nil))
}

func genRandomBytes(n int) []byte {
	b := make([]byte, n)

	if _, err := rand.Read(b); err != nil {
		log.Fatal(err)
	}

	return b
}

func isTimeoutAuthorized(timeout, maxTTL string) (bool, error) {
	t, err := tparse.AddDuration(time.Now(), timeout)
	if err != nil {
		return false, err
	}

	max, err := tparse.AddDuration(time.Now().Add(time.Second), maxTTL)
	if err != nil {
		return false, err
	}

	return t.Before(max), nil
}

package tokens

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/karrick/tparse/v2"
)

// GenRefreshToken generates a refresh token for a user in a namespace.
func GenRefreshToken(timeout, maxTTL, namespace, user string, secretKey []byte) (string, error) {
	expiry, err := authorizedExpiry(timeout, maxTTL)
	if err != nil {
		return "", err
	}

	claims := newRefreshClaims(namespace, user, expiry)
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
	expiry, err := authorizedExpiry(timeout, maxTTL)
	if err != nil {
		return "", err
	}

	claims := newAdminAccessClaims(namespaceName, userName, userId, nsId, expiry, isAdmin, isNsAdmin)
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
	expiry, err := authorizedExpiry(timeout, maxTTL)
	if err != nil {
		return "", err
	}

	claims := newAccessClaims(user, groups, orgs, expiry)
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

// TODO: this may be optimize by reusing "t" and "max" in the "expiry" computing.
func authorizedExpiry(timeout, maxTTL string) (time.Time, error) {
	t, err := tparse.AddDuration(time.Now(), timeout)
	if err != nil {
		emo.ParamError(err)
		return time.Time{}, err
	}

	max, err := tparse.AddDuration(time.Now().Add(time.Second), maxTTL)
	if err != nil {
		emo.ParamError(err)
		return time.Time{}, err
	}

	isAuthorized := t.Before(max)
	if !isAuthorized {
		emo.ParamError("Unauthorized timeout", timeout)
		return time.Time{}, nil
	}

	expiry, err := tparse.ParseNow(time.RFC3339, "now+"+timeout)
	if err != nil {
		emo.ParamError(err)
		return time.Time{}, err
	}

	return expiry.UTC(), nil
}

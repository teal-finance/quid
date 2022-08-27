package tokens

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/teal-finance/garcon/timex"
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

// NewAccessToken creates an Access Token with the JSON fields "exp", "usr", "grp" and "org".
func NewAccessToken(timeout, maxTTL, user string, groups, orgs []string, algo, secretKey string) (string, error) {
	expiry, err := authorizedExpiry(timeout, maxTTL)
	if err != nil {
		return "", err
	}

	claims := newAccessClaims(user, groups, orgs, expiry)

	method := jwt.GetSigningMethod(algo)
	if method == nil {
		err = fmt.Errorf("unsupported signing algorithm %q", algo)
		emo.ParamError(err)
		return "", err
	}
	t := jwt.NewWithClaims(method, claims)

	// convert secretKey depending on algo
	var key any
	switch method.Alg() {
	case "HS256", "HS384", "HS512":
		key = []byte(secretKey)
	case "RS256", "RS384", "RS512":
		key = &rsa.PrivateKey{ // TODO
			PublicKey:   rsa.PublicKey{},
			D:           &big.Int{},
			Primes:      []*big.Int{},
			Precomputed: rsa.PrecomputedValues{},
		}
	case "ES256", "ES384", "ES512":
		key = &ecdsa.PrivateKey{ // TODO
			PublicKey: ecdsa.PublicKey{},
			D:         &big.Int{},
		}
	case "Ed25519":
		key = ed25519.PrivateKey([]byte(secretKey))
	default:
		err = fmt.Errorf("unsupported signing algorithm %q", algo)
		emo.ParamError(err)
		return "", err
	}

	token, err := t.SignedString(key)
	if err != nil {
		emo.EncryptError(err)
		return "", err
	}

	emo.AccessToken("Issued AccessToken exp="+timeout+" usr="+user+" grp=", groups, "org=", orgs, "Algo="+algo)

	return token, nil
}

// RandomHMACKey generates a random HMAC-SHA256 key.
func RandomHMACKey() string {
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

func authorizedExpiry(timeout, maxTTL string) (time.Time, error) {
	d, err := timex.ParseDuration(timeout)
	if err != nil {
		emo.ParamError("timeout", err)
		return time.Time{}, err
	}

	max, err := timex.ParseDuration(maxTTL)
	if err != nil {
		emo.ParamError("maxTTL", err)
		return time.Time{}, err
	}

	if d > max {
		err = errors.New("Unauthorized timeout=" + timeout + " > maxTTL=" + maxTTL)
		emo.ParamError(err.Error())
		return time.Time{}, err
	}

	expiry := time.Now().Add(d).UTC()
	return expiry, nil
}

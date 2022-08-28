package tokens

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"log"
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
func NewAccessToken(timeout, maxTTL, user string, groups, orgs []string, algo string, secretKey []byte) (string, error) {
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

	key, err := convertDERToPrivateKey(algo, secretKey)
	if err != nil {
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

// convertDERToPrivateKey converts DER format to a private key depending on the algo.
func convertDERToPrivateKey(algo string, der []byte) (any, error) {
	switch algo {
	case "HS256", "HS384", "HS512":
		return der, nil
	case "RS256", "RS384", "RS512":
		return x509.ParsePKCS1PrivateKey(der)
	case "ES256", "ES384", "ES512":
		return x509.ParseECPrivateKey(der)
	case "Ed25519":
		return ed25519.PrivateKey(der), nil
	}

	err := fmt.Errorf("unsupported signing algorithm %q", algo)
	emo.ParamError(err)
	return nil, err
}

// convertDERToPrivateKey converts DER format to a private key depending on the algo.
func PrivateDERToPublicDER(algo string, der []byte) ([]byte, error) {
	switch algo {
	case "HS256", "HS384", "HS512": // HMAC: same key to sign/verify
		return der, nil
	case "RS256", "RS384", "RS512": // RSA
		private, err := x509.ParsePKCS1PrivateKey(der)
		if err != nil {
			return nil, err
		}
		public := private.Public()
		return asn1.Marshal(public)
	case "ES256", "ES384", "ES512": // ESDSA
		private, err := x509.ParseECPrivateKey(der)
		if err != nil {
			return nil, err
		}
		public := private.Public()
		return asn1.Marshal(public)
	case "Ed25519": // EdDSA
		private := ed25519.PrivateKey(der)
		public := private.Public()
		return x509.MarshalPKIXPublicKey(public)
	}

	err := fmt.Errorf("unsupported signing algorithm %q", algo)
	emo.ParamError(err)
	return nil, err
}

// GenerateHMAC generates a random HMAC-SHA256 key.
func GenerateHMAC(bits int) []byte {
	switch bits {
	case 256, 384, 512: // ok
	default:
		log.Panic("accept 256, 384 and 512 bits, but got bits=", bits)
	}

	b := genRandomBytes(bits / 8)
	h := hmac.New(sha256.New, b)
	return h.Sum(nil)
}

func genRandomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Panic(err)
	}
	return b
}

// GenerateRSA generates a random RSA private key in DER format.
func GenerateRSA(bits int) []byte {
	private, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Panicf("rsa.GenerateKey(rand.Reader, bits=%d): %v", bits, err)
	}
	return x509.MarshalPKCS1PrivateKey(private)
}

// GenerateECDSA generates a random ECDSA private key in DER format.
func GenerateECDSA(c elliptic.Curve) []byte {
	private, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		log.Panicf("ecdsa.GenerateKey(%s): %v", c.Params().Name, err)
	}

	der, err := x509.MarshalECPrivateKey(private)
	if err != nil {
		log.Panicf("GenerateECDSA(%s) x509.MarshalECPrivateKey: %v", c.Params().Name, err)
	}

	return der
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

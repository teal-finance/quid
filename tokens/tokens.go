package tokens

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/teal-finance/emo"
	"github.com/teal-finance/garcon/timex"
	"github.com/teal-finance/quid/crypt"
)

var log = emo.NewZone("tokens")

const notSupportedNotice = " not yet supported. " +
	"Please contact teal.finance@pm.me or " +
	"open an issue at https://github.com/teal-finance/quid"

// Init : init the token zone.
func Init(isVerbose, isDev, isCmd bool) {
	if !isDev && !isCmd {
		log.Verbose = emo.Yes
	}
}

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
		log.EncryptError(err)
		return "", err
	}

	log.RefreshToken("Issued a refresh token for user '" + user + "' and namespace " + namespace)
	return token, nil
}

// GenAccessToken generates an access token with HS256 signing algo.
func GenAccessToken(timeout, maxTTL, user string, groups, orgs []string, secretKey []byte) (string, error) {
	expiry, err := authorizedExpiry(timeout, maxTTL)
	if err != nil {
		return "", err
	}

	claims := newAccessClaims(user, groups, orgs, expiry)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString(secretKey)
	if err != nil {
		log.EncryptError(err)
		return "", err
	}

	log.AccessToken("Issued an access token for user", user)

	return token, nil
}

// GenAccessTokenWithAlgo creates an Access Token with the JSON fields "exp", "usr", "grp" and "org".
func GenAccessTokenWithAlgo(algo, timeout, maxTTL, user string, groups, orgs []string, secretKey []byte) (string, error) {
	expiry, err := authorizedExpiry(timeout, maxTTL)
	if err != nil {
		return "", err
	}

	claims := newAccessClaims(user, groups, orgs, expiry)

	method := jwt.GetSigningMethod(algo)
	if method == nil {
		err = fmt.Errorf("unsupported signing algorithm %q, golang-jwt supports: %+v", algo, jwt.GetAlgorithms())
		log.ParamError(err)
		return "", err
	}
	t := jwt.NewWithClaims(method, claims)

	key, err := convertDERToPrivateKey(algo, secretKey)
	if err != nil {
		return "", err
	}

	token, err := t.SignedString(key)
	if err != nil {
		log.EncryptError(err)
		return "", err
	}

	log.AccessToken("Issued AccessToken exp="+timeout+" usr="+user+" grp=", groups, "org=", orgs, "Algo="+algo)

	return token, nil
}

// convertDERToPrivateKey converts DER format to a private key depending on the algo.
func convertDERToPrivateKey(algo string, der []byte) (any, error) {
	switch algo {
	case "HS256", "HS384", "HS512":
		return der, nil
	case "RS256", "RS384", "RS512":
		return x509.ParsePKCS1PrivateKey(der)
	case "PS256", "PS384", "PS512":
		return nil, errors.New(algo + notSupportedNotice)
	case "ES256", "ES384", "ES512":
		return x509.ParseECPrivateKey(der)
	case "EdDSA":
		return ed25519.PrivateKey(der), nil
	}

	err := fmt.Errorf("unsupported signing algorithm %q, golang-jwt supports: %+v", algo, jwt.GetAlgorithms())
	log.ParamError(err)
	return nil, err
}

// PrivateToPublicDER converts a private key into a public key depending on the algo.
// The input and output are in DER form.
func PrivateToPublicDER(algo string, der []byte) ([]byte, error) {
	switch algo {
	case "HS256", "HS384", "HS512": // HMAC: same key to sign/verify
		return der, nil
	}

	public, err := PrivateToPublic(algo, der)
	if err != nil {
		return nil, err
	}

	return x509.MarshalPKIXPublicKey(public)
}

// DecryptVerificationKeyDER returns the secret key for symmetric algos (like HMAC),
// or the public key for asymmetric algos (like RSA, EdDSA).
// The returned key is in DER format.
func DecryptVerificationKeyDER(algo string, accessKey []byte) ([]byte, error) {
	private, err := crypt.AesGcmDecryptBin(accessKey)
	if err != nil {
		log.DecryptError(err)
		return nil, err
	}

	public, err := PrivateToPublicDER(algo, private)
	if err != nil {
		log.DecryptError(err)
		return nil, err
	}

	return public, nil
}

// ParsePublicDER converts a public key in DER form into the original public key.
func ParsePublicDER(algo string, der []byte) (any, error) {
	switch algo {
	case "HS256", "HS384", "HS512": // HMAC: same key to sign/verify
		return der, nil
	}

	return x509.ParsePKIXPublicKey(der)
}

// PrivateToPublic converts a private key into a public key depending on the algo.
func PrivateToPublic(algo string, der []byte) (any, error) {
	switch algo {
	case "HS256", "HS384", "HS512": // HMAC: same key to sign/verify
		return der, nil

	case "RS256", "RS384", "RS512": // RSA
		private, err := x509.ParsePKCS1PrivateKey(der)
		if err != nil {
			return nil, err
		}
		return private.Public(), nil

	case "PS256", "PS384", "PS512": // RSA + salt
		return nil, errors.New(algo + notSupportedNotice)

	case "ES256", "ES384", "ES512": // ESDSA
		private, err := x509.ParseECPrivateKey(der)
		if err != nil {
			return nil, err
		}
		return private.Public(), nil

	case "EdDSA":
		private := ed25519.PrivateKey(der)
		return private.Public(), nil
	}

	err := fmt.Errorf("unsupported signing algorithm %q, golang-jwt supports: %+v", algo, jwt.GetAlgorithms())
	log.ParamError(err)
	return nil, err
}

// GenerateSigningKeyHex produces the private key in hexadecimal form.
func GenerateSigningKeyHex(algo string) (string, error) {
	b, err := GenerateSigningKey(algo)
	return hex.EncodeToString(b), err
}

// GenerateSigningKey produces the private key of the given algorithm.
// Supported algorithms:
//
// - HS256 = HMAC using SHA-256
// - HS384 = HMAC using SHA-384
// - HS512 = HMAC using SHA-512
// - RS256 = RSASSA-PKCS-v1.5 using SHA-256
// - RS384 = RSASSA-PKCS-v1.5 using SHA-384
// - RS512 = RSASSA-PKCS-v1.5 using SHA-512
// - ES256 = ECDSA using P-256 and SHA-256
// - ES384 = ECDSA using P-384 and SHA-384
// - ES512 = ECDSA using P-521 and SHA-512
// - EdDSA = Ed25519
func GenerateSigningKey(algo string) ([]byte, error) {
	switch algo {

	// HMAC

	case "HS256":
		return GenerateKeyHMAC(256), nil
	case "HS384":
		return GenerateKeyHMAC(384), nil
	case "HS512":
		return GenerateKeyHMAC(512), nil

	// RSA: 2048 bits to prevent the error "message too long for RSA public key size"
	case "RS256", "RS384", "RS512":
		return GenerateKeyRSA(2048), nil

	// RSA + salt

	case "PS256":
		return nil, errors.New("PS256" + notSupportedNotice)
	case "PS384":
		return nil, errors.New("PS384" + notSupportedNotice)
	case "PS512":
		return nil, errors.New("PS512" + notSupportedNotice)

	// ESDSA

	case "ES256":
		return GenerateKeyECDSA(elliptic.P256()), nil
	case "ES384":
		return GenerateKeyECDSA(elliptic.P384()), nil
	case "ES512":
		return GenerateKeyECDSA(elliptic.P521()), nil

	case "EdDSA":
		return GenerateEdDSAKey(), nil
	}

	err := fmt.Errorf("unsupported signing algorithm %q, golang-jwt supports: %+v", algo, jwt.GetAlgorithms())
	log.ParamError(err)
	return nil, err
}

// GenerateKeyHMAC generates a random HMAC-SHA256 key.
func GenerateKeyHMAC(bits int) []byte {
	check := false
	if bits < 0 {
		bits = -bits
		check = true
	}

	randomBytes := GenerateRandomBytes(bits / 8)

	if check {
		var digest hash.Hash
		switch bits {
		case 256:
			digest = hmac.New(sha256.New, randomBytes)
		case 384:
			digest = hmac.New(sha512.New384, randomBytes)
		case 512:
			digest = hmac.New(sha512.New, randomBytes)
		default:
			log.Panic("accept 256, 384 and 512 bits, but got bits=", bits)
		}
		hashed := digest.Sum(nil)
		if len(randomBytes) != len(hashed) {
			log.Panicf("Unexpected key size: got=%d bytes want=%d bytes input=%d bit", len(hashed), len(randomBytes), bits)
		}
	}

	return randomBytes
}

func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Panic(err)
	}
	return b
}

// GenerateKeyRSA generates a random RSA private key in DER format.
func GenerateKeyRSA(bits int) []byte {
	private, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Panicf("rsa.GenerateKey(rand.Reader, bits=%d): %v", bits, err)
	}
	return x509.MarshalPKCS1PrivateKey(private)
}

// GenerateKeyECDSA generates a random ECDSA private key in DER format.
func GenerateKeyECDSA(c elliptic.Curve) []byte {
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

// GenerateEdDSAKey generates a random EdDSA-25519 key in DER format.
func GenerateEdDSAKey() []byte {
	_, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Panic(err)
	}

	if false {
		privateDER, err := x509.MarshalPKCS8PrivateKey(private)
		if err != nil {
			log.Panic(err)
		}
		return privateDER
	}

	return private
}

func authorizedExpiry(timeout, maxTTL string) (time.Time, error) {
	d, err := timex.ParseDuration(timeout)
	if err != nil {
		log.ParamError("timeout", err)
		return time.Time{}, err
	}

	max, err := timex.ParseDuration(maxTTL)
	if err != nil {
		log.ParamError("maxTTL", err)
		return time.Time{}, err
	}

	if d > max {
		err = errors.New("Unauthorized timeout=" + timeout + " > maxTTL=" + maxTTL)
		log.ParamError(err.Error())
		return time.Time{}, err
	}

	expiry := time.Now().Add(d).UTC()
	return expiry, nil
}

func ValidAccessToken(accessToken, algo string, verificationKeyDER []byte) error {
	verificationKey, err := ParsePublicDER(algo, verificationKeyDER)
	if err != nil {
		return err
	}

	var claims AccessClaims
	f := func(*jwt.Token) (any, error) { return verificationKey, nil }
	validator := jwt.NewParser(jwt.WithValidMethods([]string{algo}))
	token, err := validator.ParseWithClaims(accessToken, &claims, f)
	if err != nil {
		return err
	}

	if err := token.Claims.Valid(); err != nil {
		return err
	}

	return nil
}

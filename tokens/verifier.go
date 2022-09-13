package tokens

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/cristalhq/base64"
	"github.com/teal-finance/garcon/gg"
	"github.com/teal-finance/quid/server"
)

var (
	ErrThreeParts    = errors.New("JWT must be composed of three parts separated by periods")
	ErrJWTSignature  = errors.New("JWT signature mismatch")
	ErrNoBase64JWT   = errors.New("the token claims (second part of the JWT) is not base64-valid")
	ErrAlgoKeyScheme = errors.New("Unexpected AlgoKey scheme")
)

type Tokenizer interface {
	GenAccessToken(timeout, maxTTL, user string, groups, orgs []string) (string, error)
	Sign(headerPayload []byte) []byte
	Verifier
}

type Verifier interface {
	Claims(accessToken []byte) (*AccessClaims, error)
	Verify(headerPayload, signature []byte) bool
}

// NewVerifier creates a new Verifier to speed up the verification
// of many Access Tokens with the same verification key.
//
// The parameter algoKey accepts three different schemes:
//
//  1. only a HMAC secret key:
//     algoKey = "9d2e0a02121179a3c3de1b035ae1355b1548781c8ce8538a1dc0853a12dfb13d"
//
//  2. both the signing algo and its verification key:
//     algoKey = "HS256:9d2e0a02121179a3c3de1b035ae1355b1548781c8ce8538a1dc0853a12dfb13d"
//
//  3. the Quid URL to fetch the algo/key info from a given namespace
//     algoKey = "https://quid.teal.finance/v1?ns=foobar"
//
// In the two first forms, NewVerifier accepts the key to be in hexadecimal, or in Base64 form.
// NewVerifier converts the verification key into binary DER form
// depending on the key string length and the optional algo name.
// The algo name is case insensitive.
func NewVerifier(algoKey string) (Verifier, error) {
	slice := strings.SplitN(algoKey, ":", 2)
	switch len(slice) {
	case 0:
		log.Panic("NewVerifier parameter must not be empty")
	case 1:
		return NewHMAC(algoKey) // here algoKey is just the secret-key
	}

	algo := strings.ToUpper(slice[0])
	keyStr := slice[1]

	switch algo {
	case "HTTP", "HTTPS":
		return RequestAlgoKey(algoKey) // here algoKey is an URL
	case "HMAC":
		return NewHMAC(keyStr)
	case "HS256":
		return NewHS256(keyStr)
	case "HS384":
		return NewHS384(keyStr)
	case "HS512":
		return NewHS512(keyStr)
	case "RS256", "RS384", "RS512":
		return NewRSA(algo, keyStr)
	case "PS256", "PS384", "PS512":
		log.Panic(algo + notSupportedNotice)
	case "ES256":
		return NewES256(keyStr)
	case "ES384":
		return NewES384(keyStr)
	case "ES512":
		return NewES512(keyStr)
	case "EDDSA":
		return NewEdDSA(keyStr)
	}

	log.Errorf("Unexpected scheme %q in algoKey=%q", slice[0], algoKey)
	return nil, ErrAlgoKeyScheme
}

func RequestAlgoKey(uri string) (Verifier, error) {
	if p := gg.Printable(uri); p >= 0 {
		return nil, fmt.Errorf("Unprintable character at position %d in sanitized URL=%q", p, uri)
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	err = fmt.Errorf(`missing valid query parameter 'namespace' in URL: %s`, uri)

	for param, values := range u.Query() {
		switch param {
		case "ns", "namespace":
			for _, ns := range values {
				var b []byte
				b, err = server.NamespaceRequest{Namespace: ns}.MarshalJSON()
				if err != nil {
					continue
				}

				u.RawQuery = ""
				endpoint := u.String()
				var resp *http.Response
				resp, err = http.DefaultClient.Post(endpoint, "application/json", bytes.NewReader(b))
				if err != nil {
					continue
				}

				var m server.PublicKeyResponse
				err = gg.UnmarshalJSONResponse(resp, &m, 1000)
				if err != nil {
					continue
				}

				algoKey := m.Alg + ":" + m.Key
				return NewVerifier(algoKey)
			}
		}
	}

	return nil, err
}

type (
	// HS256 requires a 32-bytes secret key.
	HS256 struct{ Key []byte }
	HS384 struct{ Key []byte }
	HS512 struct{ Key []byte }
	EdDSA struct{ Key []byte }
	ES256 struct{ Key *ecdsa.PublicKey }
	ES384 struct{ Key *ecdsa.PublicKey }
	ES512 struct{ Key *ecdsa.PublicKey }
)

/*
$ go test -v ./tokens/... | grep -w Public
    tokens_test.go:155: HS256 Public  key len= 32
    tokens_test.go:155: HS512 Public  key len= 64
    tokens_test.go:155: HS384 Public  key len= 48
    tokens_test.go:155: EdDSA Public  key len= 44
    tokens_test.go:155: ES256 Public  key len= 91
    tokens_test.go:155: ES384 Public  key len= 120
    tokens_test.go:155: ES512 Public  key len= 158
    tokens_test.go:155: RS384 Public  key len= 294
    tokens_test.go:155: RS512 Public  key len= 294
    tokens_test.go:155: RS256 Public  key len= 294
*/

var (
	ErrHMACKey     = errors.New("cannot decode the HMAC key, please provide a key in hexadecimal or Base64 form (64, 96 or 128 hexadecimal digits ; 43, 64 or 86 Base64 characters)")
	ErrHS256PubKey = errors.New("cannot decode the HMAC-SHA256 key, please provide 64 hexadecimal digits or a Base64 string containing about 43 characters")
	ErrHS384PubKey = errors.New("cannot decode the HMAC-SHA384 key, please provide 96 hexadecimal digits or a Base64 string containing 64 characters")
	ErrHS512PubKey = errors.New("cannot decode the HMAC-SHA512 key, please provide 128 hexadecimal digits or a Base64 string containing about 86 characters")
	ErrEdDSAPubKey = errors.New("cannot decode the EdDSA public key, please provide 88 hexadecimal digits or a Base64 string containing about 59 characters")
	ErrES256PubKey = errors.New("cannot decode the ECDSA-P256-SHA256 public key, please provide 182 hexadecimal digits or a Base64 string containing about 122 characters")
	ErrES384PubKey = errors.New("cannot decode the ECDSA-P384-SHA384 public key, please provide 240 hexadecimal digits or a Base64 string containing 160 characters")
	ErrES512PubKey = errors.New("cannot decode the ECDSA-P512-SHA512 public key, please provide 316 hexadecimal digits or a Base64 string containing about 211 characters")
	ErrECDSAPubKey = errors.New("cannot parse the DER bytes as a valid ECDSA public key")
)

func NewHS256(keyStr string) (*HS256, error) {
	key := decodeKeyInHexOrB64(keyStr, 32)
	if key == nil {
		return nil, ErrHS256PubKey
	}
	return &HS256{key}, nil
}

func NewHS384(keyStr string) (*HS384, error) {
	key := decodeKeyInHexOrB64(keyStr, 48)
	if key == nil {
		return nil, ErrHS384PubKey
	}
	return &HS384{key}, nil
}

func NewHS512(keyStr string) (*HS512, error) {
	key := decodeKeyInHexOrB64(keyStr, 64)
	if key == nil {
		return nil, ErrHS512PubKey
	}
	return &HS512{key}, nil
}

func NewEdDSA(keyStr string) (*EdDSA, error) {
	der := decodeKeyInHexOrB64(keyStr, 44)
	if der == nil {
		return nil, ErrEdDSAPubKey
	}
	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, err
	}
	edPubKey, ok := pub.(ed25519.PublicKey)
	if !ok {
		return nil, ErrECDSAPubKey
	}
	return &EdDSA{edPubKey}, nil
}

func NewES256(keyStr string) (*ES256, error) {
	key := decodeKeyInHexOrB64(keyStr, 91)
	if key == nil {
		return nil, ErrES256PubKey
	}
	pub, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	ecPubKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrECDSAPubKey
	}
	return &ES256{ecPubKey}, nil
}

func NewES384(keyStr string) (*ES384, error) {
	key := decodeKeyInHexOrB64(keyStr, 120)
	if key == nil {
		return nil, ErrES384PubKey
	}
	pub, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	ecPubKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrECDSAPubKey
	}
	return &ES384{ecPubKey}, nil
}

func NewES512(keyStr string) (*ES512, error) {
	key := decodeKeyInHexOrB64(keyStr, 158)
	if key == nil {
		return nil, ErrES512PubKey
	}
	pub, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	ecPubKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrECDSAPubKey
	}
	return &ES512{ecPubKey}, nil
}

func NewRSA(algo, keyStr string) (Verifier, error) { return nil, nil } // TODO

func NewHMAC(keyStr string) (Tokenizer, error) {
	if tokenizer, err := NewHS256(keyStr); err == nil {
		return tokenizer, nil
	}

	if tokenizer, err := NewHS384(keyStr); err == nil {
		return tokenizer, nil
	}

	if tokenizer, err := NewHS512(keyStr); err == nil {
		return tokenizer, nil
	}

	return nil, ErrHMACKey
}

func decodeKeyInHexOrB64(keyStr string, wantLen int) (key []byte) {
	wantHex := wantLen * 2
	wantB64 := wantLen * 4 / 3

	var err error

	if len(keyStr) == wantHex {
		key, err = hex.DecodeString(keyStr)
		if err != nil {
			log.Warn(err)
			return nil
		}
	} else if wantB64-1 <= len(keyStr) && len(keyStr) <= wantB64+1 {
		key, err = base64.RawURLEncoding.DecodeString(keyStr)
		if err != nil {
			log.Warn(err)
			return nil
		}
		switch len(key) {
		case wantLen - 1:
			key = append(key, 0)
		case wantLen + 1:
			key = key[:wantLen]
		}
	} else {
		return nil
	}

	if len(key) != wantLen {
		log.Panic("want=", wantLen, "got=", len(key))
	}

	return key
}

func (v *HS256) GenAccessToken(timeout, maxTTL, user string, groups, orgs []string) (string, error) {
	return GenAccessToken(timeout, maxTTL, user, groups, orgs, v.Key)
}

func (v *HS384) GenAccessToken(timeout, maxTTL, user string, groups, orgs []string) (string, error) {
	return GenAccessTokenWithAlgo("HS384", timeout, maxTTL, user, groups, orgs, v.Key)
}

func (v *HS512) GenAccessToken(timeout, maxTTL, user string, groups, orgs []string) (string, error) {
	return GenAccessTokenWithAlgo("HS512", timeout, maxTTL, user, groups, orgs, v.Key)
}

func (v *HS256) Verify(hp, sig []byte) bool { return verify(v, hp, sig) }
func (v *HS384) Verify(hp, sig []byte) bool { return verify(v, hp, sig) }
func (v *HS512) Verify(hp, sig []byte) bool { return verify(v, hp, sig) }
func (v *ES256) Verify(hp, sig []byte) bool { return ecdsaVerify(crypto.SHA256.New(), v.Key, hp, sig) }
func (v *ES384) Verify(hp, sig []byte) bool { return ecdsaVerify(crypto.SHA384.New(), v.Key, hp, sig) }
func (v *ES512) Verify(hp, sig []byte) bool { return ecdsaVerify(crypto.SHA512.New(), v.Key, hp, sig) }
func (v *EdDSA) Verify(hp, sig []byte) bool { return ed25519.Verify(v.Key, hp, sig) }

func verify[T Tokenizer](v T, headerPayload, signature []byte) bool {
	ourSignature := v.Sign(headerPayload)
	return bytes.Equal(ourSignature, signature)
}

func ecdsaVerify(digest hash.Hash, pub *ecdsa.PublicKey, headerPayload, sig []byte) bool {
	digest.Write(headerPayload)
	r := big.NewInt(0).SetBytes(sig[:len(sig)/2])
	s := big.NewInt(0).SetBytes(sig[len(sig)/2:])
	return ecdsa.Verify(pub, digest.Sum(nil), r, s)
}

func ecdsaVerify2(digest hash.Hash, pub *ecdsa.PublicKey, headerPayload, sig []byte) bool {
	digest.Write(headerPayload)
	return ecdsa.VerifyASN1(pub, digest.Sum(nil), sig)
}

func (v *HS256) Claims(JWT []byte) (*AccessClaims, error) { return claims(v, JWT) }
func (v *HS384) Claims(JWT []byte) (*AccessClaims, error) { return claims(v, JWT) }
func (v *HS512) Claims(JWT []byte) (*AccessClaims, error) { return claims(v, JWT) }
func (v *ES256) Claims(JWT []byte) (*AccessClaims, error) { return claims(v, JWT) }
func (v *ES384) Claims(JWT []byte) (*AccessClaims, error) { return claims(v, JWT) }
func (v *ES512) Claims(JWT []byte) (*AccessClaims, error) { return claims(v, JWT) }
func (v *EdDSA) Claims(JWT []byte) (*AccessClaims, error) { return claims(v, JWT) }

func claims[T Verifier](v T, accessToken []byte) (*AccessClaims, error) {
	p1, p2, err := SplitThreeParts(accessToken)
	if err != nil {
		return nil, err
	}

	payload := accessToken[p1+1 : p2]
	ac, err := AccessClaimsFromBase64(payload)
	if err != nil {
		return nil, err
	}

	headerPayload := accessToken[:p2]
	signature := accessToken[p2+1:]
	if !v.Verify(headerPayload, signature) {
		return nil, ErrJWTSignature
	}

	return ac, nil
}

// Sign return the signature of the first two parts.
// It allocates hmac.New() each time to avoid race condition.
func (v *HS256) Sign(headerPayload []byte) []byte {
	digest := hmac.New(sha256.New, v.Key)
	digest.Write(headerPayload)
	signatureBin := digest.Sum(nil)
	signatureB64 := make([]byte, len(signatureBin)*4/3+1)
	base64.RawURLEncoding.Encode(signatureB64, signatureBin)
	return signatureB64
}

func (v *HS384) Sign(headerPayload []byte) []byte {
	digest := hmac.New(sha512.New384, v.Key)
	digest.Write(headerPayload)
	signatureBin := digest.Sum(nil)
	signatureB64 := make([]byte, len(signatureBin)*4/3+1)
	base64.RawURLEncoding.Encode(signatureB64, signatureBin)
	return signatureB64
}

func (v *HS512) Sign(headerPayload []byte) []byte {
	digest := hmac.New(sha512.New, v.Key)
	digest.Write(headerPayload)
	signatureBin := digest.Sum(nil)
	signatureB64 := make([]byte, len(signatureBin)*4/3+1)
	base64.RawURLEncoding.Encode(signatureB64, signatureBin)
	return signatureB64
}

// SplitThreeParts returns the period position decompose the JWT in three parts
func SplitThreeParts(JWT []byte) (p1, p2 int, _ error) {
	p1 = bytes.IndexByte(JWT, '.')
	p2 = bytes.LastIndexByte(JWT, '.')
	if p1 < 0 || p1 >= p2 {
		return 0, 0, ErrThreeParts
	}
	return p1, p2, nil
}

func AccessClaimsFromBase64(b64 []byte) (*AccessClaims, error) {
	claimsTxt := make([]byte, len(b64)*3/4)
	_, err := base64.RawURLEncoding.Decode(claimsTxt, b64)
	if err != nil {
		return nil, ErrNoBase64JWT
	}

	var claims AccessClaims
	if err := claims.UnmarshalJSON(claimsTxt); err != nil {
		return nil, &claimError{err, claimsTxt}
	}

	err = claims.Valid() // error can be: expired or invalid access token
	return &claims, err
}

type claimError struct {
	err         error
	claimsBytes []byte
}

func (e *claimError) Message() string {
	return "cannot JSON-decode AccessClaims"
}

func (e *claimError) Error() string {
	return e.err.Error() + " => " + e.Message() + ": " + string(e.claimsBytes)
}

func (e *claimError) Unwrap() error {
	return e.err
}

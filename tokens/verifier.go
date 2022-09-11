package tokens

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/cristalhq/base64"
)

var (
	ErrThreeParts   = errors.New("JWT must be composed of three parts separated by periods")
	ErrJWTSignature = errors.New("JWT signature mismatch")
	ErrNoBase64JWT  = errors.New("the token claims (second part of the JWT) is not base64-valid")
)

type Generator interface {
	GenAccessToken(timeout, maxTTL, user string, groups, orgs []string) (string, error)
}

type Verifier interface {
	Claims(accessToken string) (*AccessClaims, error)
	Sign(headerPayload []byte) []byte
}

type GenVerifier interface {
	Generator
	Verifier
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
func NewVerifier(algoKey string) Verifier {
	slice := strings.SplitN(algoKey, ":", 2)
	switch len(slice) {
	case 0:
		log.Panic("NewVerifier parameter must not be empty")
	case 1:
		return NewHMACVerifier(algoKey) // here: algoKey is just the secret-key
	}

	algo := strings.ToUpper(slice[0])
	keyStr := slice[1]

	switch algo {
	case "HTTP", "HTTPS":
		return RequestAlgoKey(algoKey) // here: algoKey is an URL
	case "HMAC":
		return NewHMACVerifier(keyStr)
	case "HS256":
		return NewHS256Verifier(keyStr)
	case "HS384":
		return NewHS384Verifier(keyStr)
	case "HS512":
		return NewHS512Verifier(keyStr)
	case "RS256", "RS384", "RS512":
		return NewRSAVerifier(algo, keyStr)
	case "PS256", "PS384", "PS512":
		log.Panic(algo + notSupportedNotice)
	case "ES256":
		return NewES256Verifier(keyStr)
	case "ES384":
		return NewES384Verifier(keyStr)
	case "ES512":
		return NewES512Verifier(keyStr)
	case "EDDSA":
		return NewEdDSAVerifier(keyStr)
	}

	log.Panicf("Unexpected scheme %q in algoKey=%q", slice[0], algoKey)
	return nil
}

func RequestAlgoKey(url string) Verifier            { return nil }
func NewHMACVerifier(keyStr string) GenVerifier     { return NewHS256Verifier(keyStr) }
func NewHS384Verifier(keyStr string) *HS384Verifier { return &HS384Verifier{} }
func NewHS512Verifier(keyStr string) *HS512Verifier { return &HS512Verifier{} }
func NewRSAVerifier(algo, keyStr string) Verifier   { return nil }
func NewES256Verifier(keyStr string) Verifier       { return nil }
func NewES384Verifier(keyStr string) Verifier       { return nil }
func NewES512Verifier(keyStr string) Verifier       { return nil }
func NewEdDSAVerifier(keyStr string) Verifier       { return nil }

type (
	// HS256Verifier requires a 32-bytes secret key.
	HS256Verifier struct{ Key []byte }
	HS384Verifier struct{ Key []byte }
	HS512Verifier struct{ Key []byte }
)

func (v *HS256Verifier) GenAccessToken(timeout, maxTTL, user string, groups, orgs []string) (string, error) {
	return GenAccessToken(timeout, maxTTL, user, groups, orgs, v.Key)
}

func (v *HS384Verifier) GenAccessToken(timeout, maxTTL, user string, groups, orgs []string) (string, error) {
	return GenAccessToken(timeout, maxTTL, user, groups, orgs, v.Key)
}

func (v *HS512Verifier) GenAccessToken(timeout, maxTTL, user string, groups, orgs []string) (string, error) {
	return GenAccessToken(timeout, maxTTL, user, groups, orgs, v.Key)
}

const hs256Constraint = "Cannot decode the HMAC-SHA256 key, please provide 64 hexadecimal digits or a Base64 string containing about 43 characters"

func NewHS256Verifier(secretKeyStr string) *HS256Verifier {
	if len(secretKeyStr) == 64 {
		key, err := hex.DecodeString(secretKeyStr)
		if err != nil {
			log.Panic(hs256Constraint, err)
		}
		return &HS256Verifier{key}
	}

	key, err := base64.RawURLEncoding.DecodeString(secretKeyStr)
	if err != nil {
		log.Panic(hs256Constraint, err)
	}

	switch len(key) {
	case 31:
		key = append(key, 0)
	case 32:
		// perfect
	case 33:
		key = key[:32]
	default:
		log.Panic(hs256Constraint)
	}

	return &HS256Verifier{key}
}

func (v *HS256Verifier) Claims(accessToken string) (*AccessClaims, error) {
	p1, p2, err := SplitThreeParts(accessToken)
	if err != nil {
		return nil, err
	}

	accessBytes := []byte(accessToken)
	headerPayload := accessBytes[:p2]
	signature := accessBytes[p2+1:]

	ourSignature := v.Sign(headerPayload)
	if !bytes.Equal(signature, ourSignature) {
		return nil, ErrJWTSignature
	}

	payload := accessBytes[p1+1 : p2]
	return AccessClaimsFromBase64(payload)
}

func (v *HS384Verifier) Claims(accessToken string) (*AccessClaims, error) {
	p1, p2, err := SplitThreeParts(accessToken)
	if err != nil {
		return nil, err
	}

	accessBytes := []byte(accessToken)
	headerPayload := accessBytes[:p2]
	signature := accessBytes[p2+1:]

	ourSignature := v.Sign(headerPayload)
	if !bytes.Equal(signature, ourSignature) {
		return nil, ErrJWTSignature
	}

	payload := accessBytes[p1+1 : p2]
	return AccessClaimsFromBase64(payload)
}

func (v *HS512Verifier) Claims(accessToken string) (*AccessClaims, error) {
	p1, p2, err := SplitThreeParts(accessToken)
	if err != nil {
		return nil, err
	}

	accessBytes := []byte(accessToken)
	headerPayload := accessBytes[:p2]
	signature := accessBytes[p2+1:]

	ourSignature := v.Sign(headerPayload)
	if !bytes.Equal(signature, ourSignature) {
		return nil, ErrJWTSignature
	}

	payload := accessBytes[p1+1 : p2]
	return AccessClaimsFromBase64(payload)
}

// SplitThreeParts returns the period position decompose the JWT in three parts
func SplitThreeParts(JWT string) (p1, p2 int, _ error) {
	p1 = strings.IndexByte(JWT, '.')
	p2 = strings.LastIndexByte(JWT, '.')
	if p1 < 0 || p1 >= p2 {
		return 0, 0, ErrThreeParts
	}
	return p1, p2, nil
}

// Sign return the signature of the first two parts.
// It allocates hmac.New() each time to avoid race condition.
func (v *HS256Verifier) Sign(headerPayload []byte) []byte {
	h := hmac.New(sha256.New, v.Key)
	_, _ = h.Write(headerPayload)
	signatureBin := h.Sum(nil)
	signatureB64 := make([]byte, len(signatureBin)*4/3+1)
	base64.RawURLEncoding.Encode(signatureB64, signatureBin)
	return signatureB64
}

func (v *HS384Verifier) Sign(headerPayload []byte) []byte {
	h := hmac.New(sha512.New384, v.Key)
	_, _ = h.Write(headerPayload)
	signatureBin := h.Sum(nil)
	signatureB64 := make([]byte, len(signatureBin)*4/3+1)
	base64.RawURLEncoding.Encode(signatureB64, signatureBin)
	return signatureB64
}

func (v *HS512Verifier) Sign(headerPayload []byte) []byte {
	h := hmac.New(sha512.New, v.Key)
	_, _ = h.Write(headerPayload)
	signatureBin := h.Sum(nil)
	signatureB64 := make([]byte, len(signatureBin)*4/3+1)
	base64.RawURLEncoding.Encode(signatureB64, signatureBin)
	return signatureB64
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

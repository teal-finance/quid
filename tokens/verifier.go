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
	ErrThreeParts    = errors.New("JWT must be composed of three parts separated by periods")
	ErrJWTSignature  = errors.New("JWT signature mismatch")
	ErrNoBase64JWT   = errors.New("the token claims (second part of the JWT) is not base64-valid")
	ErrHMACKey       = errors.New("cannot decode the HMAC key, please provide a key in hexadecimal or Base64 form (64, 96 or 128 hexadecimal digits ; 43, 64 or 86 Base64 characters)")
	ErrHS256Key      = errors.New("cannot decode the HMAC-SHA256 key, please provide 64 hexadecimal digits or a Base64 string containing about 43 characters")
	ErrHS384Key      = errors.New("cannot decode the HMAC-SHA384 key, please provide 96 hexadecimal digits or a Base64 string containing 64 characters")
	ErrHS512Key      = errors.New("cannot decode the HMAC-SHA512 key, please provide 128 hexadecimal digits or a Base64 string containing about 86 characters")
	ErrAlgoKeyScheme = errors.New("Unexpected AlgoKey scheme")
)

type Tokenizer interface {
	GenAccessToken(timeout, maxTTL, user string, groups, orgs []string) (string, error)
	Sign(headerPayload []byte) []byte
	Verifier
}

type Verifier interface {
	Claims(accessToken string) (*AccessClaims, error)
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

func RequestAlgoKey(url string) (Verifier, error)  { return nil, nil }
func NewRSA(algo, keyStr string) (Verifier, error) { return nil, nil }
func NewES256(keyStr string) (Verifier, error)     { return nil, nil }
func NewES384(keyStr string) (Verifier, error)     { return nil, nil }
func NewES512(keyStr string) (Verifier, error)     { return nil, nil }
func NewEdDSA(keyStr string) (Verifier, error)     { return nil, nil }

type (
	// HS256 requires a 32-bytes secret key.
	HS256 struct{ Key []byte }
	HS384 struct{ Key []byte }
	HS512 struct{ Key []byte }
)

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

func NewHS256(secretKeyStr string) (*HS256, error) {
	if len(secretKeyStr) == 64 {
		key, err := hex.DecodeString(secretKeyStr)
		if err != nil {
			log.Warn(err)
			return nil, ErrHS256Key
		}
		return &HS256{key}, nil
	}

	key, err := base64.RawURLEncoding.DecodeString(secretKeyStr)
	if err != nil {
		log.Warn(err)
		return nil, ErrHS256Key
	}

	switch len(key) {
	case 31:
		key = append(key, 0)
	case 32:
		// perfect
	case 33:
		key = key[:32]
	default:
		log.Warn("Got", len(key), "bytes, want 32")
		return nil, ErrHS256Key
	}

	return &HS256{key}, nil
}

func NewHS384(secretKeyStr string) (*HS384, error) {
	if len(secretKeyStr) == 48 {
		key, err := hex.DecodeString(secretKeyStr)
		if err != nil {
			log.Warn(err)
			return nil, ErrHS384Key
		}
		return &HS384{key}, nil
	}

	key, err := base64.RawURLEncoding.DecodeString(secretKeyStr)
	if err != nil {
		log.Warn(err)
		return nil, ErrHS384Key
	}

	if len(key) != 48 {
		log.Warn("Got", len(key), "bytes, want 48")
		return nil, ErrHS384Key
	}

	return &HS384{key}, nil
}

func NewHS512(secretKeyStr string) (*HS512, error) {
	if len(secretKeyStr) == 128 {
		key, err := hex.DecodeString(secretKeyStr)
		if err != nil {
			log.Warn(err)
			return nil, ErrHS512Key
		}
		return &HS512{key}, nil
	}

	key, err := base64.RawURLEncoding.DecodeString(secretKeyStr)
	if err != nil {
		log.Warn(err)
		return nil, ErrHS512Key
	}

	switch len(key) {
	case 63:
		key = append(key, 0)
	case 64:
		// perfect
	case 65:
		key = key[:64]
	default:
		log.Warn("Got", len(key), "bytes, want 64.")
		return nil, ErrHS512Key
	}

	return &HS512{key}, nil
}

func (v *HS256) GenAccessToken(timeout, maxTTL, user string, groups, orgs []string) (string, error) {
	return GenAccessToken(timeout, maxTTL, user, groups, orgs, v.Key)
}

func (v *HS384) GenAccessToken(timeout, maxTTL, user string, groups, orgs []string) (string, error) {
	return GenAccessToken(timeout, maxTTL, user, groups, orgs, v.Key)
}

func (v *HS512) GenAccessToken(timeout, maxTTL, user string, groups, orgs []string) (string, error) {
	return GenAccessToken(timeout, maxTTL, user, groups, orgs, v.Key)
}

func (v *HS256) Claims(JWT string) (*AccessClaims, error) { return SignedClaims(v, []byte(JWT)) }
func (v *HS384) Claims(JWT string) (*AccessClaims, error) { return SignedClaims(v, []byte(JWT)) }
func (v *HS512) Claims(JWT string) (*AccessClaims, error) { return SignedClaims(v, []byte(JWT)) }

func SignedClaims[T Tokenizer](v T, accessToken []byte) (*AccessClaims, error) {
	p1, p2, err := SplitThreeParts(accessToken)
	if err != nil {
		return nil, err
	}

	ourSignature := v.Sign(accessToken[:p2]) // pass header.payload
	if !bytes.Equal(ourSignature, accessToken[p2+1:]) {
		return nil, ErrJWTSignature
	}

	payload := accessToken[p1+1 : p2]
	return AccessClaimsFromBase64(payload)
}

// Sign return the signature of the first two parts.
// It allocates hmac.New() each time to avoid race condition.
func (v *HS256) Sign(headerPayload []byte) []byte {
	h := hmac.New(sha256.New, v.Key)
	_, _ = h.Write(headerPayload)
	signatureBin := h.Sum(nil)
	signatureB64 := make([]byte, len(signatureBin)*4/3+1)
	base64.RawURLEncoding.Encode(signatureB64, signatureBin)
	return signatureB64
}

func (v *HS384) Sign(headerPayload []byte) []byte {
	h := hmac.New(sha512.New384, v.Key)
	_, _ = h.Write(headerPayload)
	signatureBin := h.Sum(nil)
	signatureB64 := make([]byte, len(signatureBin)*4/3+1)
	base64.RawURLEncoding.Encode(signatureB64, signatureBin)
	return signatureB64
}

func (v *HS512) Sign(headerPayload []byte) []byte {
	h := hmac.New(sha512.New, v.Key)
	_, _ = h.Write(headerPayload)
	signatureBin := h.Sum(nil)
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

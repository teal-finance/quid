package tokens_test

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"reflect"
	"strings"
	"testing"

	turbobase64 "github.com/cristalhq/base64"
	"github.com/golang-jwt/jwt/v4"

	"github.com/teal-finance/emo"
	"github.com/teal-finance/quid/tokens"
)

var cases = []struct {
	name       string
	timeout    string
	maxTTL     string
	user       string
	groups     []string
	orgs       []string
	want       string
	wantGenErr bool
	wantNewErr bool
}{{
	name:       "HS256=HMAC-SHA256",
	timeout:    "11m",
	maxTTL:     "12m",
	user:       "me",
	groups:     []string{"dev"},
	orgs:       []string{"wikipedia"},
	want:       "",
	wantGenErr: false,
	wantNewErr: false,
}, {
	name:       "HS384=HMAC-SHA384",
	timeout:    "11m",
	maxTTL:     "12m",
	user:       "me",
	groups:     []string{"dev"},
	orgs:       []string{"wikipedia"},
	want:       "",
	wantGenErr: false,
	wantNewErr: false,
}, {
	name:       "HS512=HMAC-SHA512",
	timeout:    "11m",
	maxTTL:     "12m",
	user:       "me",
	groups:     []string{"dev"},
	orgs:       []string{"wikipedia"},
	want:       "",
	wantGenErr: false,
	wantNewErr: false,
}, {
	name:       "RS256=RSASSA-PKCSv15-SHA256",
	timeout:    "11m",
	maxTTL:     "12m",
	user:       "",
	groups:     nil,
	orgs:       nil,
	want:       "",
	wantGenErr: false,
	wantNewErr: false,
}, {
	name:       "RS384=RSASSA-PKCSv15-SHA384",
	timeout:    "11m",
	maxTTL:     "12m",
	user:       "",
	groups:     nil,
	orgs:       nil,
	want:       "",
	wantGenErr: false,
	wantNewErr: false,
}, {
	name:       "RS512=RSASSA-PKCSv15-SHA512",
	timeout:    "11m",
	maxTTL:     "12m",
	user:       "",
	groups:     nil,
	orgs:       nil,
	want:       "",
	wantGenErr: false,
	wantNewErr: false,
}, {
	name:       "ES256=ECDSA-P256-SHA256",
	timeout:    "11m",
	maxTTL:     "12m",
	user:       "me",
	groups:     []string{"dev"},
	orgs:       []string{"wikipedia"},
	want:       "",
	wantGenErr: false,
	wantNewErr: false,
}, {
	name:       "ES384=ECDSA-P384-SHA384",
	timeout:    "11m",
	maxTTL:     "12m",
	user:       "me",
	groups:     []string{"dev"},
	orgs:       []string{"wikipedia"},
	want:       "",
	wantGenErr: false,
	wantNewErr: false,
}, {
	name:       "ES512=ECDSA-P521-SHA512",
	timeout:    "11m",
	maxTTL:     "12m",
	user:       "me",
	groups:     []string{"dev"},
	orgs:       []string{"wikipedia"},
	want:       "",
	wantGenErr: false,
	wantNewErr: false,
}, {
	name:       "EdDSA=Ed25519",
	timeout:    "11m",
	maxTTL:     "12m",
	user:       "me",
	groups:     []string{"dev"},
	orgs:       []string{"wikipedia"},
	want:       "",
	wantGenErr: false,
	wantNewErr: false,
}}

func TestNewAccessToken(t *testing.T) {
	t.Parallel()

	emo.GlobalColoring(false)

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			algo := strings.Split(c.name, "=")[0]

			privateKey, err := tokens.GenerateSigningKey(algo)
			if (err != nil) != c.wantGenErr {
				t.Errorf("GenerateSigningKey() error = %v, wantGenErr %v", err, c.wantGenErr)
				return
			}
			t.Log(algo+" Private key len=", len(privateKey))

			publicDER, err := tokens.PrivateToPublicDER(algo, privateKey)
			if err != nil {
				t.Error("Public(privateKey) error:", err)
				return
			}
			t.Log(algo+" Public  key len=", len(publicDER))

			publicKey, err := tokens.PrivateToPublic(algo, privateKey)
			if err != nil {
				t.Error("PrivateToPublic("+algo+",privateKey) error:", err)
				return
			}

			publicKey2, err := tokens.ParsePublicDER(algo, publicDER)
			if err != nil {
				t.Error("ParsePublicDER("+algo+",der) error:", err)
				return
			}

			if !reflect.DeepEqual(publicKey2, publicKey) {
				t.Error("public keys are not equal")
			}

			tokenStr, err := tokens.GenAccessTokenWithAlgo(c.timeout, c.maxTTL, c.user, c.groups, c.orgs, algo, privateKey)
			if (err != nil) != c.wantNewErr {
				t.Errorf("NewAccessToken() error = %v, wantNewErr %v", err, c.wantNewErr)
				return
			}
			t.Log(algo+" TokenString len=", len(tokenStr), tokenStr)

			validator := jwt.NewParser(jwt.WithValidMethods([]string{algo}))

			var claims tokens.AccessClaims
			f := func(*jwt.Token) (any, error) { return publicKey, nil }
			token, err := validator.ParseWithClaims(tokenStr, &claims, f)
			if err != nil {
				t.Error("ParseWithClaims error:", err)
				return
			}

			if err := token.Claims.Valid(); err != nil {
				t.Error("token.Claims.Valid:", err)
			}

			if err := tokens.ValidAccessToken(tokenStr, algo, publicDER); err != nil {
				t.Error("ValidAccessToken:", err)
			}
		})
	}
}

/* To verify the JWT signature, what is the most pertinent?
   1. sign header+payload, then base64.Encode our signature, and finally compare the two base64 signatures
   2. sign header+payload, then base64.Decode their signature, and finally compare the two binary signatures
   Are the strings much less performant than []byte?

   Let's bench it:

go test -bench=. -benchmem ./...
goos: linux
goarch: amd64
pkg: github.com/teal-finance/quid/tokens
cpu: AMD Ryzen 9 3900X 12-Core Processor
BenchmarkBase64Encode-24         5913382   209.2 ns/op   144 B/op   1 allocs/op
BenchmarkBase64EncodeTurbo-24    8053846   144.6 ns/op   144 B/op   1 allocs/op
BenchmarkBase64EncodeString-24   3928515   297.4 ns/op   256 B/op   2 allocs/op
BenchmarkBase64Decode-24         4375579   277.0 ns/op   112 B/op   1 allocs/op
BenchmarkBase64DecodeTurbo-24    5554828   206.1 ns/op   112 B/op   1 allocs/op
BenchmarkBase64DecodeString-24   3331362   376.8 ns/op   256 B/op   2 allocs/op
PASS
ok      github.com/teal-finance/quid/tokens     6.254s
*/

func BenchmarkBase64Encode(b *testing.B) {
	srcTxt := []byte(`{"usr": "jane","grp": ["group1", "group2"],"org": ["organization1", "organization2"],"exp": 1595950745}`)
	b64BytesSame := make([]byte, len(srcTxt)*4/3+1)
	base64.RawURLEncoding.Encode(b64BytesSame, []byte(srcTxt))
	b64BytesDiff := make([]byte, len(srcTxt)*4/3+1)
	_, _ = rand.Read(b64BytesDiff)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dstBase64Bytes := make([]byte, len(srcTxt)*4/3+1)
		base64.RawURLEncoding.Encode(dstBase64Bytes, srcTxt)
		_ = dstBase64Bytes

		var ok bool
		if n%2 == 0 {
			ok = bytes.Equal(dstBase64Bytes, b64BytesSame)
		} else {
			ok = bytes.Equal(dstBase64Bytes, b64BytesDiff)
		}
		_ = ok
	}
}

func BenchmarkBase64EncodeTurbo(b *testing.B) {
	srcTxt := []byte(`{"usr": "jane","grp": ["group1", "group2"],"org": ["organization1", "organization2"],"exp": 1595950745}`)
	b64BytesSame := make([]byte, len(srcTxt)*4/3+1)
	base64.RawURLEncoding.Encode(b64BytesSame, []byte(srcTxt))
	b64BytesDiff := make([]byte, len(srcTxt)*4/3+1)
	_, _ = rand.Read(b64BytesDiff)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dstBase64Bytes := make([]byte, len(srcTxt)*4/3+1)
		turbobase64.RawURLEncoding.Encode(dstBase64Bytes, srcTxt)
		_ = dstBase64Bytes

		var ok bool
		if n%2 == 0 {
			ok = bytes.Equal(dstBase64Bytes, b64BytesSame)
		} else {
			ok = bytes.Equal(dstBase64Bytes, b64BytesDiff)
		}
		_ = ok
	}
}

func BenchmarkBase64EncodeString(b *testing.B) {
	srcTxt := `{"usr": "jane","grp": ["group1", "group2"],"org": ["organization1", "organization2"],"exp": 1595950745}`
	b64BytesSame := make([]byte, len(srcTxt)*4/3+1)
	base64.RawURLEncoding.Encode(b64BytesSame, []byte(srcTxt))
	b64BytesDiff := make([]byte, len(srcTxt)*4/3+1)
	_, _ = rand.Read(b64BytesDiff)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dstBase64Bytes := make([]byte, len(srcTxt)*4/3+1)
		base64.RawURLEncoding.Encode(dstBase64Bytes, []byte(srcTxt))

		var ok bool
		if n%2 == 0 {
			ok = bytes.Equal(dstBase64Bytes, b64BytesSame)
		} else {
			ok = bytes.Equal(dstBase64Bytes, b64BytesDiff)
		}
		_ = ok
	}
}

func BenchmarkBase64Decode(b *testing.B) {
	txtSame := []byte(`{"usr": "jane","grp": ["group1", "group2"],"org": ["organization1", "organization2"],"exp": 1595950745}`)
	txtDiff := make([]byte, len(txtSame))
	_, _ = rand.Read(txtDiff)

	b64 := make([]byte, len(txtSame)*4/3+1)
	base64.RawURLEncoding.Encode(b64, txtSame)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dstTxt := make([]byte, len(txtSame))
		_, err := base64.RawURLEncoding.Decode(dstTxt, b64)
		if err != nil {
			panic(err)
		}

		var ok bool
		if n%2 == 0 {
			ok = bytes.Equal(dstTxt, txtSame)
		} else {
			ok = bytes.Equal(dstTxt, txtDiff)
		}
		_ = ok
	}
}

func BenchmarkBase64DecodeTurbo(b *testing.B) {
	txtSame := []byte(`{"usr": "jane","grp": ["group1", "group2"],"org": ["organization1", "organization2"],"exp": 1595950745}`)
	txtDiff := make([]byte, len(txtSame))
	_, _ = rand.Read(txtDiff)

	b64 := make([]byte, len(txtSame)*4/3+1)
	base64.RawURLEncoding.Encode(b64, txtSame)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dstTxt := make([]byte, len(txtSame))
		_, err := turbobase64.RawURLEncoding.Decode(dstTxt, b64)
		if err != nil {
			panic(err)
		}

		var ok bool
		if n%2 == 0 {
			ok = bytes.Equal(dstTxt, txtSame)
		} else {
			ok = bytes.Equal(dstTxt, txtDiff)
		}
		_ = ok
	}
}

func BenchmarkBase64DecodeString(b *testing.B) {
	txtSame := []byte(`{"usr": "jane","grp": ["group1", "group2"],"org": ["organization1", "organization2"],"exp": 1595950745}`)
	txtDiff := make([]byte, len(txtSame))
	_, _ = rand.Read(txtDiff)

	b64Bytes := make([]byte, len(txtSame)*4/3+1)
	base64.RawURLEncoding.Encode(b64Bytes, txtSame)
	b64String := string(b64Bytes)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dstTxt, err := base64.RawURLEncoding.DecodeString(b64String)
		if err != nil {
			panic(err)
		}

		var ok bool
		if n%2 == 0 {
			ok = bytes.Equal(dstTxt, txtSame)
		} else {
			ok = bytes.Equal(dstTxt, txtDiff)
		}
		_ = ok
	}
}

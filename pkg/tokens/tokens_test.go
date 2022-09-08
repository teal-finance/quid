package tokens_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/teal-finance/emo"
	"github.com/teal-finance/quid/pkg/tokens"
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

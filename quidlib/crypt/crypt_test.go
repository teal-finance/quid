package crypt

import (
	"testing"

	"github.com/teal-finance/quid/quidlib/conf"
)

func TestAesGcm(t *testing.T) {
	conf.EncodingKey = []byte("eb037d66a3d07cc90c393a9bb04c172c")

	data := "some plaintext"
	out, err := AesGcmEncryptHex(data)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	in, err := AesGcmDecryptHex(out)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if data != in {
		t.Fatalf("expect %x got %x", data, in)
	}
}

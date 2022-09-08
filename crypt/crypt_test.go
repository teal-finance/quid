package crypt_test

import (
	"testing"

	"github.com/teal-finance/quid/crypt"
)

func TestAesGcm(t *testing.T) {
	crypt.EncodingKey = []byte("eb037d66a3d07cc90c393a9bb04c172c")

	data := "some plaintext"
	out, err := crypt.AesGcmEncryptHex(data)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	in, err := crypt.AesGcmDecryptHex(out)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if data != in {
		t.Fatalf("expect %x got %x", data, in)
	}
}

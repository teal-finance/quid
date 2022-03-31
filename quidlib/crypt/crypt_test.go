package crypt

import (
	"testing"

	"github.com/teal-finance/quid/quidlib/conf"
)

func TestAesGcm(t *testing.T) {
	conf.EncodingKey = "eb037d66a3d07cc90c393a9bb04c172c"
	data := "someplaintext"
	out, err := AesGcmEncrypt(data, nil)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}
	in, err := AesGcmDecrypt(out, nil)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}
	if data != in {
		t.Fatalf("expect %x got %x", data, in)
	}
}

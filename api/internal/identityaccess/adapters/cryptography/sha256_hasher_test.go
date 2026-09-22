package cryptography

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

func TestSHA256HasherReturnsSHA256Digest(t *testing.T) {
	value := "token-secreto"
	want := sha256.Sum256([]byte(value))

	got, err := (SHA256Hasher{}).Hash(value)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want[:]) {
		t.Fatalf("unexpected hash: %x", got)
	}
}

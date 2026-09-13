package cryptography

import (
	"encoding/base64"
	"testing"
)

func TestRandomTokenGeneratorGeneratesURLSafe32ByteTokens(t *testing.T) {
	generator := RandomTokenGenerator{}
	first, err := generator.GenerateLoginToken()
	if err != nil {
		t.Fatal(err)
	}
	second, err := generator.GenerateLoginToken()
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("expected independently generated tokens to differ")
	}

	raw, err := base64.RawURLEncoding.DecodeString(first)
	if err != nil {
		t.Fatalf("expected URL-safe base64 token: %v", err)
	}
	if len(raw) != 32 {
		t.Fatalf("expected 32 random bytes, got %d", len(raw))
	}
}

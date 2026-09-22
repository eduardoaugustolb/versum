package cryptography

import (
	"crypto/hmac"
	"crypto/sha256"
)

func HMACSHA256(message, key []byte) []byte {
	hash := hmac.New(sha256.New, key)
	_, _ = hash.Write(message)
	return hash.Sum(nil)
}

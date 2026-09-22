package keyring

import (
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/config"
)

func TestKeyRingStaticReturnsPreviousEncryptionKey(t *testing.T) {
	keyring := NewKeyRingStatic(config.Config{
		EnvironmentVariable: config.EnvironmentVariable{
			EncryptionKeyVersion: 2,
			EncryptionKey:        []byte("01234567890123456789012345678901"),
			EncryptionKeys: map[int][]byte{
				1: []byte("abcdefghijklmnop"),
				2: []byte("01234567890123456789012345678901"),
			},
		},
	})

	key, err := keyring.EncryptionKey(1)
	if err != nil {
		t.Fatal(err)
	}
	if string(key) != "abcdefghijklmnop" {
		t.Fatalf("unexpected previous key: %q", key)
	}
}

func TestKeyRingStaticReturnsCurrentEncryptionKeyWithoutKeyMap(t *testing.T) {
	keyring := NewKeyRingStatic(config.Config{
		EnvironmentVariable: config.EnvironmentVariable{
			EncryptionKeyVersion: 1,
			EncryptionKey:        []byte("01234567890123456789012345678901"),
		},
	})

	key, err := keyring.EncryptionKey(1)
	if err != nil {
		t.Fatal(err)
	}
	if string(key) != "01234567890123456789012345678901" {
		t.Fatalf("unexpected current key: %q", key)
	}
}

func TestKeyRingStaticReturnsVersionedLookupKeys(t *testing.T) {
	keyring := NewKeyRingStatic(config.Config{
		EnvironmentVariable: config.EnvironmentVariable{
			LookupKeyVersion: 2,
			LookupKey:        []byte("abcdefghijklmnopqrstuvwxyz123456"),
			LookupKeys: map[int][]byte{
				1: []byte("abcdefghijklmnop"),
				2: []byte("abcdefghijklmnopqrstuvwxyz123456"),
			},
		},
	})

	previousKey, err := keyring.LookupKey(1)
	if err != nil {
		t.Fatal(err)
	}
	if string(previousKey) != "abcdefghijklmnop" {
		t.Fatalf("unexpected previous key: %q", previousKey)
	}

	keys, err := keyring.LookupKeys()
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 || string(keys[2]) != "abcdefghijklmnopqrstuvwxyz123456" {
		t.Fatalf("unexpected lookup keys: %#v", keys)
	}
}

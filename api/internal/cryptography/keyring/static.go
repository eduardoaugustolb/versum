package keyring

import (
	"errors"
	"maps"

	"github.com/eduardoaugustolb/versum/api/internal/config"
)

type KeyRingStatic struct {
	cfg config.Config
}

var _ Keyring = (*KeyRingStatic)(nil)

func NewKeyRingStatic(config config.Config) *KeyRingStatic {
	return &KeyRingStatic{
		cfg: config,
	}
}

func (k *KeyRingStatic) CurrentEncryptionKey() (version int, key []byte, err error) {
	return k.cfg.EncryptionKeyVersion, k.cfg.EncryptionKey, nil
}

func (k *KeyRingStatic) EncryptionKey(version int) ([]byte, error) {
	key, ok := k.cfg.EncryptionKeys[version]
	if ok {
		return key, nil
	}

	if version == k.cfg.EncryptionKeyVersion {
		return k.cfg.EncryptionKey, nil
	}

	return nil, errors.New("invalid encryption key version")
}

func (k *KeyRingStatic) CurrentLookupKey() (version int, key []byte, err error) {
	return k.cfg.LookupKeyVersion, k.cfg.LookupKey, nil
}

func (k *KeyRingStatic) LookupKey(version int) ([]byte, error) {
	key, ok := k.cfg.LookupKeys[version]
	if ok {
		return key, nil
	}

	if version == k.cfg.LookupKeyVersion {
		return k.cfg.LookupKey, nil
	}
	return nil, errors.New("invalid lookup key version")
}

func (k *KeyRingStatic) LookupKeys() (map[int][]byte, error) {
	keys := make(map[int][]byte, len(k.cfg.LookupKeys)+1)
	maps.Copy(keys, k.cfg.LookupKeys)
	if _, ok := keys[k.cfg.LookupKeyVersion]; !ok {
		keys[k.cfg.LookupKeyVersion] = k.cfg.LookupKey
	}
	return keys, nil
}

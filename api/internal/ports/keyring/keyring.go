package keyring

type Keyring interface {
	CurrentEncryptionKey() (version int, key []byte, err error)
	EncryptionKey(version int) ([]byte, error)

	CurrentLookupKey() (version int, key []byte, err error)
	LookupKey(version int) ([]byte, error)
}

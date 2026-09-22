package config

import "errors"

var (
	ErrPortNotNumeric = errors.New("port must be a number")
	ErrPortOutOfRange = errors.New("port must be between 1 and 65535")

	ErrInvalidEnvironment = errors.New("invalid environment")

	ErrDatabaseURLNotSet             = errors.New("database url not set")
	ErrEncryptionSecretNotSet        = errors.New("encryption secret not set")
	ErrLookupSecretNotSet            = errors.New("lookup secret not set")
	ErrInvalidEncryptionKey          = errors.New("invalid encryption key: must be 16, 24 or 32 bytes")
	ErrInvalidEncryptionVersion      = errors.New("invalid encryption key version: must be a positive integer")
	ErrInvalidPreviousEncryptionKeys = errors.New("invalid previous encryption keys")
	ErrInvalidLookupKey              = errors.New("invalid lookup key: must be at least 16 bytes")
	ErrInvalidLookupKeyVersion       = errors.New("invalid lookup key version: must be a positive integer")
	ErrInvalidPreviousLookupKeys     = errors.New("invalid previous lookup keys")

	// Deprecated: use ErrLookupSecretNotSet.
	ErrLookupSecretVersionNotSet = ErrLookupSecretNotSet
	// Deprecated: use ErrInvalidEncryptionVersion.
	ErrInvalidSecretVersion = ErrInvalidEncryptionVersion
)

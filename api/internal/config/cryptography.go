package config

const (
	DefaultEncryptionSecretKey     = "ENCRYPTION_SECRET_KEY"
	DefaultEncryptionSecretVersion = "ENCRYPTION_SECRET_VERSION"
	DefaultPreviousEncryptionKeys  = "ENCRYPTION_SECRET_PREVIOUS_KEYS"
	DefaultLookupSecretKey         = "LOOKUP_SECRET_KEY"
	DefaultLookupSecretVersion     = "LOOKUP_SECRET_VERSION"
	DefaultPreviousLookupKeys      = "LOOKUP_SECRET_PREVIOUS_KEYS"
	MaxLookupKeyVersions           = 32767
	MaxEncryptionSecretKeyVersions = 32767
)

// Deprecated: use DefaultLookupSecretKey.
const DefualtLookupSecretKey = DefaultLookupSecretKey

// Deprecated: use DefaultLookupSecretVersion.
const DefualtLookupSecretVersion = DefaultLookupSecretVersion

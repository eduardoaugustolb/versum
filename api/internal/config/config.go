package config

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Config struct {
	Address     string
	Environment Environment
	EnvironmentVariable
}

type EnvironmentVariable struct {
	DatabaseURL          string
	RedisURL             string
	LookupKey            []byte
	LookupKeyVersion     int
	LookupKeys           map[int][]byte
	EncryptionKey        []byte
	EncryptionKeyVersion int
	EncryptionKeys       map[int][]byte
}

func Load(lookup func(string) string) (Config, error) {
	cfg := Config{}

	port, err := parsePort(lookup("PORT"))
	if err != nil {
		return Config{}, err
	}
	cfg.Address = ":" + port

	// Environment
	environment, err := parseEnvironment(lookup("ENVIRONMENT"))
	if err != nil {
		return Config{}, err
	}
	cfg.Environment = environment

	// Database
	databaseURL := lookup(DefaultDatabaseURLKey)
	if databaseURL == "" {
		return Config{}, ErrDatabaseURLNotSet
	}
	cfg.DatabaseURL = databaseURL

	// Redis
	redisURL := lookup(DefaultRedisURLKey)
	if redisURL == "" {
		return Config{}, ErrRedisURLNotSet
	}
	cfg.RedisURL = redisURL

	// Encryption
	encryptionKey := lookup(DefaultEncryptionSecretKey)
	encryptionKeyVersionRaw := lookup(DefaultEncryptionSecretVersion)

	if encryptionKey == "" || encryptionKeyVersionRaw == "" {
		return Config{}, ErrEncryptionSecretNotSet
	}
	// len([]byte) counts bytes, not runes, so multibyte chars are accounted
	// for. AES-GCM only accepts 16, 24 or 32 byte keys.
	switch len([]byte(encryptionKey)) {
	case 16, 24, 32:
	default:
		return Config{}, ErrInvalidEncryptionKey
	}
	encryptionKeyVersion, err := strconv.Atoi(encryptionKeyVersionRaw)
	if err != nil {
		return Config{}, fmt.Errorf("%w: %w", ErrInvalidEncryptionVersion, err)
	}
	if encryptionKeyVersion <= 0 {
		return Config{}, ErrInvalidEncryptionVersion
	}
	if encryptionKeyVersion > MaxEncryptionSecretKeyVersions {
		return Config{}, ErrInvalidEncryptionKey
	}

	cfg.EncryptionKey = []byte(encryptionKey)
	cfg.EncryptionKeyVersion = encryptionKeyVersion
	cfg.EncryptionKeys = map[int][]byte{encryptionKeyVersion: cfg.EncryptionKey}

	if rawPreviousKeys := lookup(DefaultPreviousEncryptionKeys); rawPreviousKeys != "" {
		var previousKeys map[string]string
		if err := json.Unmarshal([]byte(rawPreviousKeys), &previousKeys); err != nil {
			return Config{}, fmt.Errorf("%w: %w", ErrInvalidPreviousEncryptionKeys, err)
		}
		for rawVersion, previousKey := range previousKeys {
			version, err := strconv.Atoi(rawVersion)
			if err != nil || version <= 0 || version == encryptionKeyVersion || version > MaxLookupKeyVersions {
				return Config{}, ErrInvalidPreviousEncryptionKeys
			}
			switch len([]byte(previousKey)) {
			case 16, 24, 32:
			default:
				return Config{}, ErrInvalidPreviousEncryptionKeys
			}
			cfg.EncryptionKeys[version] = []byte(previousKey)
		}
	}

	// Lookup
	lookupKey := lookup(DefaultLookupSecretKey)
	lookupKeyVersionRaw := lookup(DefaultLookupSecretVersion)
	if lookupKey == "" || lookupKeyVersionRaw == "" {
		return Config{}, ErrLookupSecretNotSet
	}

	if len([]byte(lookupKey)) < 16 {
		return Config{}, ErrInvalidLookupKey
	}
	lookupKeyVersion, err := strconv.Atoi(lookupKeyVersionRaw)
	if err != nil {
		return Config{}, fmt.Errorf("%w: %w", ErrInvalidLookupKeyVersion, err)
	}
	if lookupKeyVersion <= 0 {
		return Config{}, ErrInvalidLookupKeyVersion
	}
	if lookupKeyVersion > MaxLookupKeyVersions {
		return Config{}, ErrInvalidLookupKeyVersion
	}

	cfg.LookupKey = []byte(lookupKey)
	cfg.LookupKeyVersion = lookupKeyVersion
	cfg.LookupKeys = map[int][]byte{lookupKeyVersion: cfg.LookupKey}

	if rawPreviousKeys := lookup(DefaultPreviousLookupKeys); rawPreviousKeys != "" {
		var previousKeys map[string]string
		if err := json.Unmarshal([]byte(rawPreviousKeys), &previousKeys); err != nil {
			return Config{}, fmt.Errorf("%w: %w", ErrInvalidPreviousLookupKeys, err)
		}
		for rawVersion, previousKey := range previousKeys {
			version, err := strconv.Atoi(rawVersion)
			if err != nil || version <= 0 || version == lookupKeyVersion || len([]byte(previousKey)) < 16 || version > MaxLookupKeyVersions {
				return Config{}, ErrInvalidPreviousLookupKeys
			}
			cfg.LookupKeys[version] = []byte(previousKey)
		}
	}

	return cfg, nil
}

package config_test

import (
	"errors"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/config"
)

func testLookup(overrides map[string]string) func(string) string {
	values := map[string]string{
		"ENVIRONMENT":                         string(config.DefaultEnvironment),
		"PORT":                                config.DefaultPort,
		config.DefaultDatabaseURLKey:          config.DefaultDatabaseURL,
		config.DefaultRedisURLKey:             "redis://127.0.0.1:6379/0",
		config.DefaultEncryptionSecretKey:     "01234567890123456789012345678901",
		config.DefaultEncryptionSecretVersion: "1",
		config.DefaultPreviousEncryptionKeys:  "",
		config.DefaultLookupSecretKey:         "abcdefghijklmnopqrstuvwxyz123456",
		config.DefaultLookupSecretVersion:     "1",
		config.DefaultPreviousLookupKeys:      "",
	}
	for k, v := range overrides {
		values[k] = v
	}
	return func(key string) string {
		return values[key]
	}
}

func TestLoadPort(t *testing.T) {
	tests := []struct {
		name        string
		port        string
		wantErr     error
		wantAddress string
	}{
		{
			name:    "not numeric port",
			port:    "808A",
			wantErr: config.ErrPortNotNumeric,
		},
		{
			name:    "port above range",
			port:    "70000",
			wantErr: config.ErrPortOutOfRange,
		},
		{
			name:    "negative port",
			port:    "-1",
			wantErr: config.ErrPortOutOfRange,
		},
		{
			name:        "port unset uses default",
			port:        "",
			wantAddress: ":8080",
		},
		{
			name:        "valid port",
			port:        "9090",
			wantAddress: ":9090",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lookup := testLookup(map[string]string{
				"PORT": tc.port,
			})

			cfg, err := config.Load(lookup)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if cfg.Address != tc.wantAddress {
				t.Errorf("expected address %q, got %q", tc.wantAddress, cfg.Address)
			}
		})
	}
}

func TestLoadEnvironment(t *testing.T) {
	tests := []struct {
		name            string
		environment     string
		wantErr         error
		wantEnvironment config.Environment
	}{
		{
			name:            "environment unset uses default",
			environment:     "",
			wantEnvironment: config.DevelopmentEnvironment,
		},
		{
			name:            "development",
			environment:     "development",
			wantEnvironment: config.DevelopmentEnvironment,
		},
		{
			name:            "production",
			environment:     "production",
			wantEnvironment: config.ProductionEnvironment,
		},
		{
			name:        "invalid environment",
			environment: "staging",
			wantErr:     config.ErrInvalidEnvironment,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lookup := testLookup(map[string]string{
				"ENVIRONMENT": tc.environment,
			})

			cfg, err := config.Load(lookup)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if cfg.Environment != tc.wantEnvironment {
				t.Errorf("expected environment %q, got %q", tc.wantEnvironment, cfg.Environment)
			}
		})
	}
}

func TestLoadDatabaseURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr error
	}{
		{
			name:    "database URL unset returns error",
			url:     "",
			wantErr: config.ErrDatabaseURLNotSet,
		},
		{
			name: "database URL set",
			url:  config.DefaultDatabaseURL,
			want: config.DefaultDatabaseURL,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lookup := testLookup(map[string]string{
				config.DefaultDatabaseURLKey: tc.url,
			})

			cfg, err := config.Load(lookup)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if cfg.DatabaseURL != tc.want {
				t.Errorf("expected database URL %q, got %q", tc.want, cfg.DatabaseURL)
			}
		})
	}
}

func TestLoadSecrets(t *testing.T) {
	tests := []struct {
		name      string
		overrides map[string]string
		wantErr   error
	}{
		{
			name: "encryption secret unset",
			overrides: map[string]string{
				config.DefaultEncryptionSecretKey: "",
			},
			wantErr: config.ErrEncryptionSecretNotSet,
		},
		{
			name: "encryption key short",
			overrides: map[string]string{
				config.DefaultEncryptionSecretKey: "short",
			},
			wantErr: config.ErrInvalidEncryptionKey,
		},
		{
			name: "encryption version not numeric",
			overrides: map[string]string{
				config.DefaultEncryptionSecretVersion: "abc",
			},
			wantErr: config.ErrInvalidEncryptionVersion,
		},
		{
			name: "encryption version zero",
			overrides: map[string]string{
				config.DefaultEncryptionSecretVersion: "0",
			},
			wantErr: config.ErrInvalidEncryptionVersion,
		},
		{
			name: "lookup secret unset",
			overrides: map[string]string{
				config.DefaultLookupSecretKey: "",
			},
			wantErr: config.ErrLookupSecretNotSet,
		},
		{
			name: "lookup key short",
			overrides: map[string]string{
				config.DefaultLookupSecretKey: "short",
			},
			wantErr: config.ErrInvalidLookupKey,
		},
		{
			name: "lookup version not numeric",
			overrides: map[string]string{
				config.DefaultLookupSecretVersion: "abc",
			},
			wantErr: config.ErrInvalidLookupKeyVersion,
		},
		{
			name: "lookup version zero",
			overrides: map[string]string{
				config.DefaultLookupSecretVersion: "0",
			},
			wantErr: config.ErrInvalidLookupKeyVersion,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lookup := testLookup(tc.overrides)

			_, err := config.Load(lookup)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestLoadPreviousEncryptionKeys(t *testing.T) {
	cfg, err := config.Load(testLookup(map[string]string{
		config.DefaultEncryptionSecretVersion: "2",
		config.DefaultPreviousEncryptionKeys:  `{"1":"abcdefghijklmnop"}`,
	}))
	if err != nil {
		t.Fatal(err)
	}

	if string(cfg.EncryptionKeys[1]) != "abcdefghijklmnop" {
		t.Fatalf("unexpected previous key: %q", cfg.EncryptionKeys[1])
	}
	if string(cfg.EncryptionKeys[2]) != "01234567890123456789012345678901" {
		t.Fatalf("unexpected current key: %q", cfg.EncryptionKeys[2])
	}
}

func TestLoadRejectsInvalidPreviousEncryptionKeys(t *testing.T) {
	_, err := config.Load(testLookup(map[string]string{
		config.DefaultPreviousEncryptionKeys: `{"invalid":"abcdefghijklmnop"}`,
	}))
	if !errors.Is(err, config.ErrInvalidPreviousEncryptionKeys) {
		t.Fatalf("expected invalid previous encryption keys, got %v", err)
	}
}

func TestLoadPreviousLookupKeys(t *testing.T) {
	cfg, err := config.Load(testLookup(map[string]string{
		config.DefaultLookupSecretVersion: "2",
		config.DefaultPreviousLookupKeys:  `{"1":"abcdefghijklmnop"}`,
	}))
	if err != nil {
		t.Fatal(err)
	}

	if string(cfg.LookupKeys[1]) != "abcdefghijklmnop" {
		t.Fatalf("unexpected previous key: %q", cfg.LookupKeys[1])
	}
	if string(cfg.LookupKeys[2]) != "abcdefghijklmnopqrstuvwxyz123456" {
		t.Fatalf("unexpected current key: %q", cfg.LookupKeys[2])
	}
}

func TestLoadRejectsInvalidPreviousLookupKeys(t *testing.T) {
	_, err := config.Load(testLookup(map[string]string{
		config.DefaultPreviousLookupKeys: `{"invalid":"abcdefghijklmnop"}`,
	}))
	if !errors.Is(err, config.ErrInvalidPreviousLookupKeys) {
		t.Fatalf("expected invalid previous lookup keys, got %v", err)
	}
}

package config_test

import (
	"errors"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/config"
)

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
			lookup := func(key string) string {
				return map[string]string{
					"ENVIRONMENT": string(config.DefaultEnvironment),
					"PORT":        tc.port,
				}[key]
			}

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
			lookup := func(key string) string {
				return map[string]string{
					"ENVIRONMENT": tc.environment,
					"PORT":        config.DefaultPort,
				}[key]
			}

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

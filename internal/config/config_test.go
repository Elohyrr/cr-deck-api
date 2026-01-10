package config

import (
	"os"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"SUPERCELL_API_KEY": "test_api_key",
				"API_TOKEN":         "test_bearer_token",
				"POSTGRES_PASSWORD": "test_password",
				"TOP_PLAYERS_LIMIT": "500",
				"API_PORT":          "9090",
			},
			wantErr: false,
		},
		{
			name: "missing SUPERCELL_API_KEY",
			envVars: map[string]string{
				"API_TOKEN":         "test_token",
				"POSTGRES_PASSWORD": "test_password",
			},
			wantErr: true,
		},
		{
			name: "missing API_TOKEN",
			envVars: map[string]string{
				"SUPERCELL_API_KEY": "test_key",
				"POSTGRES_PASSWORD": "test_password",
			},
			wantErr: true,
		},
		{
			name: "missing POSTGRES_PASSWORD",
			envVars: map[string]string{
				"SUPERCELL_API_KEY": "test_key",
				"API_TOKEN":         "test_token",
			},
			wantErr: true,
		},
		{
			name: "invalid TOP_PLAYERS_LIMIT (too high)",
			envVars: map[string]string{
				"SUPERCELL_API_KEY": "test_key",
				"API_TOKEN":         "test_token",
				"POSTGRES_PASSWORD": "test_password",
				"TOP_PLAYERS_LIMIT": "1500",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg, err := LoadFromEnv()

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if cfg.SupercellAPIKey != tt.envVars["SUPERCELL_API_KEY"] {
					t.Errorf("SupercellAPIKey = %v, want %v", cfg.SupercellAPIKey, tt.envVars["SUPERCELL_API_KEY"])
				}
				if cfg.APIToken != tt.envVars["API_TOKEN"] {
					t.Errorf("APIToken = %v, want %v", cfg.APIToken, tt.envVars["API_TOKEN"])
				}
			}
		})
	}
}

func TestConfig_PostgresDSN(t *testing.T) {
	cfg := &Config{
		PostgresHost:     "localhost",
		PostgresPort:     5432,
		PostgresUser:     "testuser",
		PostgresPassword: "testpass",
		PostgresDB:       "testdb",
	}

	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	got := cfg.PostgresDSN()

	if got != expected {
		t.Errorf("PostgresDSN() = %v, want %v", got, expected)
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				SupercellAPIKey:  "key",
				APIToken:         "token",
				PostgresPassword: "pass",
				TopPlayersLimit:  500,
			},
			wantErr: false,
		},
		{
			name: "missing SupercellAPIKey",
			config: &Config{
				APIToken:         "token",
				PostgresPassword: "pass",
				TopPlayersLimit:  500,
			},
			wantErr: true,
		},
		{
			name: "TopPlayersLimit too high",
			config: &Config{
				SupercellAPIKey:  "key",
				APIToken:         "token",
				PostgresPassword: "pass",
				TopPlayersLimit:  2000,
			},
			wantErr: true,
		},
		{
			name: "TopPlayersLimit too low",
			config: &Config{
				SupercellAPIKey:  "key",
				APIToken:         "token",
				PostgresPassword: "pass",
				TopPlayersLimit:  0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

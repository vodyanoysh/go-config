package goconfig

import (
	"os"
	"path/filepath"
	"testing"
)

type TestConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Features []string       `yaml:"features"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

func TestLoadConfig(t *testing.T) {
	// Создаем временный файл конфигурации
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `
database:
  host: localhost
  port: 5432
  user: [[ getenv "DB_USER" ]]
  password: [[ getenv "DB_PASSWORD" ]]
server:
  address: [[ getenv "SERVER_ADDRESS" ]]
  port: 8080
features:
  - feature1
  - [[ getenv "FEATURE_NAME" ]]
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("failed to create test config file: %v", err)
	}

	// Устанавливаем переменные окружения
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("SERVER_ADDRESS", "0.0.0.0")
	os.Setenv("FEATURE_NAME", "dynamic_feature")
	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("SERVER_ADDRESS")
		os.Unsetenv("FEATURE_NAME")
	}()

	// Загружаем конфигурацию
	var config TestConfig
	err = LoadConfig(&config, configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Проверяем значения
	if config.Database.User != "testuser" {
		t.Errorf("expected DB user 'testuser', got '%s'", config.Database.User)
	}
	if config.Database.Password != "testpass" {
		t.Errorf("expected DB password 'testpass', got '%s'", config.Database.Password)
	}
	if config.Server.Address != "0.0.0.0" {
		t.Errorf("expected server address '0.0.0.0', got '%s'", config.Server.Address)
	}
	if len(config.Features) != 2 {
		t.Errorf("expected 2 features, got %d", len(config.Features))
	}
	if len(config.Features) > 1 && config.Features[1] != "dynamic_feature" {
		t.Errorf("expected feature 'dynamic_feature', got '%s'", config.Features[1])
	}
}

func TestReplaceEnvVarInString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		envVar   string
		envValue string
		expected string
	}{
		{
			name:     "simple replacement",
			input:    `[[ getenv "TEST_VAR" ]]`,
			envVar:   "TEST_VAR",
			envValue: "test_value",
			expected: "test_value",
		},
		{
			name:     "replacement in middle",
			input:    `prefix [[ getenv "TEST_VAR" ]] suffix`,
			envVar:   "TEST_VAR",
			envValue: "middle",
			expected: "prefix middle suffix",
		},
		{
			name:     "multiple replacements",
			input:    `[[ getenv "VAR1" ]] and [[ getenv "VAR2" ]]`,
			envVar:   "VAR1",
			envValue: "value1",
			expected: "value1 and value2",
		},
		{
			name:     "no replacement",
			input:    `plain text`,
			envVar:   "",
			envValue: "",
			expected: "plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVar != "" {
				os.Setenv(tt.envVar, tt.envValue)
				defer os.Unsetenv(tt.envVar)
			}

			// Для теста с множественными заменами
			if tt.name == "multiple replacements" {
				os.Setenv("VAR2", "value2")
				defer os.Unsetenv("VAR2")
			}

			result := replaceEnvVarInString(tt.input)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

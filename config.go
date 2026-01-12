package goconfig

import (
	"fmt"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var envVarPattern = regexp.MustCompile(`\[\[\s*getenv\s+"([^"]+)"\s*\]\]`)

// LoadConfig загружает конфигурацию из YAML файла и заменяет переменные окружения.
// Поддерживает формат [[ getenv "ENV_NAME" ]] для подстановки переменных окружения.
// Автоматически загружает .env файл, если он присутствует.
func LoadConfig(config any, path string) error {
	_ = godotenv.Load()

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading the file: %w", err)
	}

	// Заменяем переменные окружения в тексте YAML файла перед парсингом
	yamlContent := string(data)
	yamlContent = replaceEnvVarInString(yamlContent)

	err = yaml.Unmarshal([]byte(yamlContent), config)
	if err != nil {
		return fmt.Errorf("error unmarshalling the file: %w", err)
	}

	return nil
}

// replaceEnvVarInString заменяет все вхождения [[ getenv "ENV_NAME" ]] на значения переменных окружения
func replaceEnvVarInString(s string) string {
	return envVarPattern.ReplaceAllStringFunc(s, func(match string) string {
		matches := envVarPattern.FindStringSubmatch(match)
		if len(matches) > 1 {
			envVar := matches[1]
			envValue := os.Getenv(envVar)
			if envValue != "" {
				return envValue
			}
			// Если переменная не найдена, оставляем оригинальную строку
			return match
		}
		return match
	})
}

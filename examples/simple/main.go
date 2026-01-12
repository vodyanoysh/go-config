package main

import (
	"log"
	"os"

	goconfig "github.com/vodyanoysh/go-config"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	App      AppConfig      `yaml:"app"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type AppConfig struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Debug       bool   `yaml:"debug"`
}

func main() {
	// Устанавливаем примеры переменных окружения
	os.Setenv("DB_USER", "admin")
	os.Setenv("DB_PASSWORD", "secret123")
	os.Setenv("SERVER_ADDRESS", "0.0.0.0")
	os.Setenv("APP_ENV", "production")

	var config Config

	err := goconfig.LoadConfig(&config, "examples/simple/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Configuration loaded successfully!")
	log.Printf("Database Host: %s", config.Database.Host)
	log.Printf("Database User: %s", config.Database.User)
	log.Printf("Server Address: %s", config.Server.Address)
	log.Printf("Server Port: %d", config.Server.Port)
	log.Printf("App Name: %s", config.App.Name)
	log.Printf("App Environment: %s", config.App.Environment)
	log.Printf("Debug Mode: %v", config.App.Debug)
}

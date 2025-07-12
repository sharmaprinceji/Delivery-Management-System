package config

import (
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

type Config struct {
	Env         string     `yaml:"env" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
}

var (
	cfg  *Config
	once sync.Once
)

func MustLoad() *Config {
	once.Do(func() {
		// Hardcoded fallback config path
		const configPath = "config/local.yaml"

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("Configuration file does not exist: %s", configPath)
		}

		var c Config
		if err := cleanenv.ReadConfig(configPath, &c); err != nil {
			log.Fatalf("Failed to read config: %v", err)
		}

		cfg = &c
	})

	return cfg
}

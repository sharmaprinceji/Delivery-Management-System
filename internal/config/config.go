package config


import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

type Config struct {
	Env         string      `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string      `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer  `yaml:"http_server"`
}

var (
	cfg  *Config
	once sync.Once
)

func MustLoad() *Config {
	once.Do(func() {
		// var configPath string

		// Priority 1: CONFIG_PATH environment variable
		var configPath = os.Getenv("CONFIG_PATH")

		// Priority 2: --config flag
		if configPath == "" {
			flag.StringVar(&configPath, "config", "", "Path to configuration file")
			flag.Parse()
		}

		// Error if config path not found
		if configPath == "" {
			log.Fatal("CONFIG_PATH environment variable or --config flag must be set")
		}

		// Validate file existence
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
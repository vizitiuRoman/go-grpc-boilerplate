package config

import (
	"os"
	"time"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/db"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/server/grpc"
	"go.uber.org/config"
	"go.uber.org/zap"
)

type Config struct {
	Location string `yaml:"location"`

	Name      string `yaml:"name"`
	Version   string `yaml:"version"`
	BuildDate string `yaml:"build_date"`

	Logger *zap.Config  `yaml:"logger"`
	DB     *db.Config   `yaml:"db"`
	Server *grpc.Config `yaml:"server"`
}

var cfg *Config

const defaultConfigYaml = "./config.yaml"

var (
	name      = "saga-svc"
	version   = "undefined/local"
	buildDate = time.Now().Format(time.RFC3339)
)

func NewConfig() (*Config, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = defaultConfigYaml
	}

	provider, err := NewProviderByOptions(config.File(configFile))
	if err != nil {
		return nil, err
	}
	if err = provider.Populate(&cfg); err != nil {
		panic(err)
	}

	cfg.Name = name
	cfg.Version = version
	cfg.BuildDate = buildDate

	return cfg, nil
}

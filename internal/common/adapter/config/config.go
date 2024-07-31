package config

import (
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/server/grpc"
	"github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/config"
	log "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/logger"
	"github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/pgclient"
)

type Config struct {
	Location string `yaml:"location"`

	Name      string `yaml:"name"`
	Version   string `yaml:"version"`
	BuildDate string `yaml:"build_date"`

	Logger *log.Config      `yaml:"logger"`
	DB     *pgclient.Config `yaml:"pgclient"`
	Server *grpc.Config     `yaml:"server"`
}

var cfg *Config

func NewConfig() (*Config, error) {
	return cfg, config.PopulateConfig(cfg)
}

package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type configOptions struct {
	logger         Logger
	yamlConfigPath string
	envFilePath    string
}

type Option interface {
	apply(*configOptions)
}

type optionFunc func(*configOptions)

func (f optionFunc) apply(c *configOptions) { f(c) }

func WithLogger(logger Logger) Option {
	return optionFunc(func(c *configOptions) {
		c.logger = logger
	})
}

func WithConfig(path string) Option {
	return optionFunc(func(c *configOptions) {
		c.yamlConfigPath = path
	})
}

func WithENV(path string) Option {
	return optionFunc(func(c *configOptions) {
		c.envFilePath = path
	})
}

func (o *configOptions) getLogger() Logger {
	if o.logger == nil {
		o.logger = zap.NewExample()
	}

	return o.logger
}

func readENV(path string) (map[string]string, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	dotenvFile, err := godotenv.Read(path)
	if err != nil {
		return nil, err
	}

	return dotenvFile, nil
}

func (o *configOptions) localENV() (map[string]string, error) {
	if o.envFilePath != "" {
		if res, err := readENV(o.envFilePath); err != nil {
			return nil, err
		} else {
			return res, nil
		}
	}

	for _, path := range []string{os.Getenv("ENV_FILE"), "./.env"} {
		if res, err := readENV(path); err == nil {
			o.getLogger().Info("using env file", zap.String("path", path))
			return res, nil
		}
	}

	return nil, nil
}

func (o *configOptions) getYamlConfigPath() string {
	if o.yamlConfigPath != "" {
		return o.yamlConfigPath
	}

	configFile := os.Getenv("CONFIG_FILE")
	if _, err := readENV(configFile); err == nil {
		o.getLogger().Info("using config file", zap.String("path", configFile))
		return configFile
	}

	return "./config.yaml"
}

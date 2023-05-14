package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/config"
)

type Provider interface {
	Populate(any) error
	PopulateByKey(string, any) error
}

const defaultEnv = "./.env"

func NewProviderByOptions(options ...config.YAMLOption) (Provider, error) {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = defaultEnv
	}

	dotenvFile, err := godotenv.Read(envFile)
	if err != nil {
		return nil, err
	}

	configProvider, err := config.NewYAML(append(
		[]config.YAMLOption{
			config.Expand(lookupFunc(dotenvFile)),
		},
		options...,
	)...)
	if err != nil {
		return nil, err
	}

	return &ymlProvider{configProvider}, nil
}

func lookupFunc(extraEnv map[string]string) func(key string) (val string, ok bool) {
	return func(key string) (val string, ok bool) {
		val, ok = os.LookupEnv(key)
		if !ok {
			val, ok = extraEnv[key]
		}
		return val, ok
	}
}

type ymlProvider struct {
	provider config.Provider
}

func (p *ymlProvider) Populate(target any) error {
	return p.PopulateByKey("", target)
}

func (p *ymlProvider) PopulateByKey(key string, target any) error {
	return p.provider.Get(key).Populate(target)
}

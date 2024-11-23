package config

import (
	"os"

	core "go.uber.org/config"
)

type Provider interface {
	Populate(any) error
	PopulateByKey(string, any) error
}

func newProvider(options *configOptions) (Provider, error) {
	dotenvFile, err := options.localENV()
	if err != nil {
		return nil, err
	}

	configProvider, err := core.NewYAML(
		core.File(options.getYamlConfigPath()),
		core.Expand(lookupFunc(dotenvFile)))
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
	provider core.Provider
}

func (p *ymlProvider) Populate(target any) error {
	return p.PopulateByKey("", target)
}

func (p *ymlProvider) PopulateByKey(key string, target any) error {
	return p.provider.Get(key).Populate(target)
}

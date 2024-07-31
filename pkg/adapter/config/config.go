package config

func PopulateConfig(cfg any, opts ...Option) error {
	o := &configOptions{}
	for _, opt := range opts {
		opt.apply(o)
	}

	provider, err := newProvider(o)
	if err != nil {
		return err
	}

	return provider.Populate(cfg)
}

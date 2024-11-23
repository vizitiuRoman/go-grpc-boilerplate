package pgclient

type Option func(cfg *pgDB)

func WithHookFactories(hookFactories []HookFactory) Option {
	return func(cln *pgDB) {
		cln.hookFactories = hookFactories
	}
}

func WithHookFactory(hookFactory HookFactory) Option {
	return func(cln *pgDB) {
		cln.hookFactories = append(cln.hookFactories, hookFactory)
	}
}

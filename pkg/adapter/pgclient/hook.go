package pgclient

import (
	"context"
)

type Hook interface {
	Before(context.Context)
	After(ctx context.Context, query string, args ...any)
}

type HookFactory interface {
	Create() Hook
}

type hookFactories []HookFactory

func (f hookFactories) createHooks() []Hook {
	if len(f) == 0 {
		return nil
	}

	hooks := make([]Hook, len(f))

	for i := 0; i < len(f); i++ {
		hooks[i] = f[i].Create()
	}

	return hooks
}

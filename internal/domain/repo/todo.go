package repo

import (
	"context"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/model"
	"github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/pgclient"
)

type TodoRepoFactory interface {
	Create(context.Context, pgclient.DB) TodoRepo
}

type TodoRepo interface {
	Find(context.Context, int64) (*model.Todo, error)
	FindAll(ctx context.Context) ([]*model.Todo, error)
	Create(context.Context, *model.Todo) (*model.Todo, error)
	Update(context.Context, *model.Todo) (*model.Todo, error)
	Delete(context.Context, int64) error
}

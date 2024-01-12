package service

import (
	"context"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/model"
	pb "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/todo/v1"
)

type TodoService interface {
	Find(context.Context, int64) (*model.Todo, error)
	FindAll(ctx context.Context) ([]*model.Todo, error)
	Create(context.Context, *pb.CreateTodoInput) (*model.Todo, error)
	Update(context.Context, *pb.UpdateTodoInput) (*model.Todo, error)
	Delete(context.Context, int64) error
}

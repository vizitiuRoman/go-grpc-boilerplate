package service

import (
	"context"

	"github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/api/pb"
)

type TodoService interface {
	Find(context.Context, int64) (*pb.Todo, error)
	FindAll(ctx context.Context) ([]*pb.Todo, error)
	Create(context.Context, *pb.CreateTodoInput) (*pb.Todo, error)
	Update(context.Context, *pb.UpdateTodoInput) (*pb.Todo, error)
	Delete(context.Context, int64) error
}

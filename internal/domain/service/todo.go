package service

import (
	"context"

	pb "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/todo/v1"
)

type TodoService interface {
	Find(context.Context, int64) (*pb.Todo, error)
	FindAll(ctx context.Context) ([]*pb.Todo, error)
	Create(context.Context, *pb.CreateTodoInput) (*pb.Todo, error)
	Update(context.Context, *pb.UpdateTodoInput) (*pb.Todo, error)
	Delete(context.Context, int64) error
}

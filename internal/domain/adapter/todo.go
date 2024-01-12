package adapter

import (
	"context"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/model"
	"github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/sqlboiler/tododb"
	pb "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/todo/v1"
)

type TodoAdapterFactory interface {
	Create(context.Context) TodoAdapter
}

type TodoAdapter interface {
	FromProto(todo *pb.Todo) *model.Todo
	ToProto(todo *model.Todo) *pb.Todo
	ToProtos(todos []*model.Todo) []*pb.Todo

	ToEntity(*model.Todo) *tododb.Todo
	ToEntities([]*model.Todo) tododb.TodoSlice
	FromEntity(*tododb.Todo) *model.Todo
	FromEntities(tododb.TodoSlice) []*model.Todo
}

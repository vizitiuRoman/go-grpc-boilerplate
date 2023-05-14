package adapter

import (
	"context"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/model"
	"github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/api/pb"
)

type TodoAdapterFactory interface {
	Create(context.Context) TodoAdapter
}

type TodoAdapter interface {
	FromProto(todo *pb.Todo) *model.Todo
	ToProto(todo *model.Todo) *pb.Todo
	ToProtos(todos []*model.Todo) []*pb.Todo
}

package adapter

import (
	"context"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/log"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/adapter"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/model"
	"github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/api/pb"
)

type todoAdapterFactory struct {
	logger log.Logger
}

func NewTodoAdapterFactory(logger log.Logger) adapter.TodoAdapterFactory {
	return &todoAdapterFactory{
		logger: logger,
	}
}

func (f *todoAdapterFactory) Create(ctx context.Context) adapter.TodoAdapter {
	return &todoAdapter{
		logger: f.logger.WithComponent(ctx, "ActionConditionAdapter"),
	}
}

type todoAdapter struct {
	logger log.Logger
}

func (t todoAdapter) FromProto(todo *pb.Todo) *model.Todo {
	return &model.Todo{
		ID:          todo.Id,
		Name:        todo.Name,
		Description: todo.Description,
	}
}

func (t todoAdapter) ToProto(todo *model.Todo) *pb.Todo {
	return &pb.Todo{
		Id:          todo.ID,
		Name:        todo.Name,
		Description: todo.Description,
	}
}

func (t todoAdapter) ToProtos(todos []*model.Todo) []*pb.Todo {
	output := make([]*pb.Todo, 0, len(todos))

	for _, todo := range todos {
		output = append(output, &pb.Todo{
			Id:          todo.ID,
			Name:        todo.Name,
			Description: todo.Description,
		})
	}

	return output
}

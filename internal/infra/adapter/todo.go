package adapter

import (
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/model"
	log "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/logger"
	"github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/sqlboiler/tododb"
	pb "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/todo/v1"
)

type TodoAdapter struct {
	logger log.Logger
}

func NewTodoAdapter(logger log.Logger) *TodoAdapter {
	return &TodoAdapter{
		logger: logger,
	}
}

func (t *TodoAdapter) FromProto(todo *pb.Todo) *model.Todo {
	return &model.Todo{
		ID:          todo.Id,
		Name:        todo.Name,
		Description: todo.Description,
	}
}

func (t *TodoAdapter) ToProto(todo *model.Todo) *pb.Todo {
	return &pb.Todo{
		Id:          todo.ID,
		Name:        todo.Name,
		Description: todo.Description,
	}
}

func (t *TodoAdapter) ToProtos(todos []*model.Todo) []*pb.Todo {
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

func (t *TodoAdapter) ToEntity(todo *model.Todo) *tododb.Todo {
	return &tododb.Todo{
		ID:          int(todo.ID),
		Name:        todo.Name,
		Description: todo.Description,
	}
}

func (t *TodoAdapter) ToEntities(todos []*model.Todo) tododb.TodoSlice {
	entities := make(tododb.TodoSlice, 0, len(todos))

	for _, todo := range todos {
		entities = append(entities, t.ToEntity(todo))
	}

	return entities
}

func (t *TodoAdapter) FromEntity(todo *tododb.Todo) *model.Todo {
	return &model.Todo{
		ID:          int64(todo.ID),
		Name:        todo.Name,
		Description: todo.Description,
	}
}

func (t *TodoAdapter) FromEntities(todos tododb.TodoSlice) []*model.Todo {
	todosModel := make([]*model.Todo, 0, len(todos))

	for _, todo := range todos {
		todosModel = append(todosModel, t.FromEntity(todo))
	}

	return todosModel
}

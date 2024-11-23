package service

import (
	"context"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/db"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/log"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/adapter"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/model"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/repo"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/service"
	pb "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/todo/v1"
	"go.uber.org/zap"
)

type todoService struct {
	logger log.Logger
	db     db.DB

	todoAdapterFactory adapter.TodoAdapterFactory
	todoRepoFactory    repo.TodoRepoFactory
}

func NewTodoService(
	ctx context.Context,
	logger log.Logger,
	db db.DB,

	todoAdapterFactory adapter.TodoAdapterFactory,
	todoRepoFactory repo.TodoRepoFactory,
) service.TodoService {
	return &todoService{
		db: db,

		todoAdapterFactory: todoAdapterFactory,
		todoRepoFactory:    todoRepoFactory,

		logger: logger.WithComponent(ctx, "todo_service"),
	}
}

func (s *todoService) Find(ctx context.Context, id int64) (*model.Todo, error) {
	todoRepo := s.todoRepoFactory.Create(ctx, s.db)

	todo, err := todoRepo.Find(ctx, id)
	if err != nil {
		s.logger.
			WithMethod(ctx, "Find").
			Error(
				"cannot find todo by id",
				zap.Int64("id", id),
				zap.Error(err),
			)
		return nil, err
	}

	return todo, nil
}

func (s *todoService) FindAll(ctx context.Context) ([]*model.Todo, error) {
	todoRepo := s.todoRepoFactory.Create(ctx, s.db)

	todos, err := todoRepo.FindAll(ctx)
	if err != nil {
		s.logger.
			WithMethod(ctx, "FindAll").
			Error(
				"cannot find all todos",
				zap.Error(err),
			)
		return nil, err
	}

	return todos, nil
}

func (s *todoService) Create(ctx context.Context, input *pb.CreateTodoInput) (*model.Todo, error) {
	todoAdapter := s.todoAdapterFactory.Create(ctx)
	todoRepo := s.todoRepoFactory.Create(ctx, s.db)

	todo, err := todoRepo.Create(ctx, todoAdapter.FromProto(input.Todo))
	if err != nil {
		s.logger.
			WithMethod(ctx, "Create").
			Error(
				"cannot create a todo",
				zap.Error(err),
			)
		return nil, err
	}

	return todo, nil
}

func (s *todoService) Update(ctx context.Context, input *pb.UpdateTodoInput) (*model.Todo, error) {
	todoAdapter := s.todoAdapterFactory.Create(ctx)
	todoRepo := s.todoRepoFactory.Create(ctx, s.db)

	todo, err := todoRepo.Update(ctx, todoAdapter.FromProto(input.Todo))
	if err != nil {
		s.logger.
			WithMethod(ctx, "Update").
			Error(
				"cannot update the todo",
				zap.Error(err),
			)
		return nil, err
	}

	return todo, nil
}

func (s *todoService) Delete(ctx context.Context, id int64) error {
	todoRepo := s.todoRepoFactory.Create(ctx, s.db)

	err := todoRepo.Delete(ctx, id)
	if err != nil {
		s.logger.
			WithMethod(ctx, "Delete").
			Error(
				"cannot delete the todo",
				zap.Error(err),
			)

		return err
	}

	return nil
}

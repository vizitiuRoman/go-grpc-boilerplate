package service

import (
	"context"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/model"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/infra/adapter"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/infra/repo"
	log "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/logger"
	"github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/pgclient"
	pb "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/todo/v1"
	"go.uber.org/zap"
)

type TodoService struct {
	logger log.Logger
	db     pgclient.DB

	adapter         *adapter.TodoAdapter
	todoRepoFactory *repo.TodoRepoFactory
}

func NewTodoService(
	ctx context.Context,
	logger log.Logger,
	db pgclient.DB,

	adapter *adapter.TodoAdapter,
	todoRepoFactory *repo.TodoRepoFactory,
) *TodoService {
	return &TodoService{
		db: db,

		adapter:         adapter,
		todoRepoFactory: todoRepoFactory,

		logger: logger.WithComponent(ctx, "todo_service"),
	}
}

func (s *TodoService) Find(ctx context.Context, id int64) (*model.Todo, error) {
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

func (s *TodoService) FindAll(ctx context.Context) ([]*model.Todo, error) {
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

func (s *TodoService) Create(ctx context.Context, input *pb.CreateTodoInput) (*model.Todo, error) {
	todoRepo := s.todoRepoFactory.Create(ctx, s.db)

	todo, err := todoRepo.Create(ctx, s.adapter.FromProto(input.Todo))
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

// Update --> @TODO --> return error when todo is created with the same ID
func (s *TodoService) Update(ctx context.Context, input *pb.UpdateTodoInput) (*model.Todo, error) {
	todoRepo := s.todoRepoFactory.Create(ctx, s.db)

	todo, err := todoRepo.Update(ctx, s.adapter.FromProto(input.Todo))
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

func (s *TodoService) Delete(ctx context.Context, id int64) error {
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

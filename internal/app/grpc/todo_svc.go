package grpc

import (
	"context"
	"errors"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/server/grpc"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/infra/adapter"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/infra/service"
	log "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/logger"
	"github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/pgclient"
	pb "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/todo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type todoSVC struct {
	pb.UnimplementedTodoSVCServer
	ctx    context.Context
	logger log.Logger
	db     pgclient.DB

	adapter *adapter.TodoAdapter

	todoService service.TodoService
}

func NewTodoSVCServerDescriptor(
	ctx context.Context,
	logger log.Logger,
	db pgclient.DB,

	adapter *adapter.TodoAdapter,

	todoService service.TodoService,
) *grpc.ServerDescriptor {
	server := &todoSVC{
		ctx: ctx,
		db:  db,

		adapter: adapter,

		todoService: todoService,

		logger: logger,
	}

	return &grpc.ServerDescriptor{
		Server:               server,
		GRPCRegistrar:        pb.RegisterTodoSVCServer,
		GRPCGatewayRegistrar: pb.RegisterTodoSVCHandlerFromEndpoint,
		Methods: []grpc.MethodDescriptor{
			{
				Method: (*todoSVC).CreateTodo,
			},
			{
				Method: (*todoSVC).UpdateTodo,
			},
			{
				Method: (*todoSVC).DeleteTodo,
			},
			{
				Method: (*todoSVC).GetTodo,
			},
			{
				Method: (*todoSVC).GetTodos,
			},
		},
	}
}

func (s *todoSVC) CreateTodo(ctx context.Context, input *pb.CreateTodoInput) (*pb.CreateTodoOutput, error) {
	todo, err := s.todoService.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTodoOutput{Todo: s.adapter.ToProto(todo)}, nil
}

func (s *todoSVC) UpdateTodo(ctx context.Context, input *pb.UpdateTodoInput) (*pb.UpdateTodoOutput, error) {
	todo, err := s.todoService.Update(ctx, input)

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return nil, status.Error(codes.NotFound, err.Error())

	case errors.Is(err, domain.ErrAlreadyExists):
		return nil, status.Error(codes.AlreadyExists, err.Error())

	case err != nil:
		return nil, err

	default:
		return &pb.UpdateTodoOutput{Todo: s.adapter.ToProto(todo)}, nil
	}
}

func (s *todoSVC) DeleteTodo(ctx context.Context, input *pb.DeleteTodoInput) (*emptypb.Empty, error) {
	err := s.todoService.Delete(ctx, input.Id)

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return nil, status.Error(codes.NotFound, err.Error())

	case err != nil:
		return nil, err

	default:
		return &emptypb.Empty{}, nil
	}
}

func (s *todoSVC) GetTodo(ctx context.Context, input *pb.GetTodoInput) (*pb.GetTodoOutput, error) {
	todo, err := s.todoService.Find(ctx, input.Id)

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return nil, status.Error(codes.NotFound, err.Error())

	case err != nil:
		return nil, err

	default:
		return &pb.GetTodoOutput{Todo: s.adapter.ToProto(todo)}, nil
	}
}

func (s *todoSVC) GetTodos(ctx context.Context, _ *pb.GetTodosInput) (*pb.GetTodosOutput, error) {
	todos, err := s.todoService.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetTodosOutput{Todos: s.adapter.ToProtos(todos)}, nil
}

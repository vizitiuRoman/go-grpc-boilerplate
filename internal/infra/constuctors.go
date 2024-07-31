package infra

import (
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/infra/adapter"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/infra/repo"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/infra/service"
	"go.uber.org/fx"
)

var Constructors = fx.Provide(
	repo.NewTodoRepoFactory,

	adapter.NewTodoAdapter,

	service.NewTodoService,
)

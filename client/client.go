package main

import (
	"context"
	"fmt"

	log "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/logger"
	todov1 "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/todo/v1"
	pb "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/todo/v1"
)

func main() {
	logger, err := log.NewLogger(nil)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	client, err := todov1.NewTodoSVCClient(
		ctx,
		logger.WithComponent(ctx, "todo-svc-client"),
		&todov1.Config{
			Host:     "localhost",
			GrpcPort: 3774,
		},
		nil, nil,
	)

	todos, err := client.GetTodos(context.Background(), &pb.GetTodosInput{})
	if err != nil {
		panic(err)
	}

	fmt.Println(todos.Todos[0])
}

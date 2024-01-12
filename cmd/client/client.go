package main

import (
	"context"
	"fmt"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/log"
	todov1 "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/todo/v1"
	pb "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/todo/v1"
)

func main() {
	logger := log.MustDefaultConsoleLogger("debug")

	client, err := todov1.NewTodoSVCClient(
		context.Background(),
		logger,
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

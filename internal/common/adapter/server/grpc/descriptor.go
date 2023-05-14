package grpc

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type ServerDescriptor struct {
	Server               any
	Methods              []MethodDescriptor
	GRPCRegistrar        any
	GRPCGatewayRegistrar func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
}

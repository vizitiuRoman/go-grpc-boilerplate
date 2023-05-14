package grpc

import (
	"context"

	"google.golang.org/grpc"
)

type ContextKey string

const MethodDescriptorContextKey = ContextKey("MethodDescriptor")

func GetDescriptor(ctx context.Context) *MethodDescriptor {
	v := ctx.Value(MethodDescriptorContextKey)
	if v == nil {
		return nil
	}
	md, ok := v.(*MethodDescriptor)
	if !ok {
		return nil
	}
	return md
}

func SetDescriptor(ctx context.Context, methodDesc *MethodDescriptor) context.Context {
	return context.WithValue(ctx, MethodDescriptorContextKey, methodDesc)
}

func MakeMethodDescriptorInterceptor(methodDescriptors MethodDescriptors) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md := methodDescriptors.GetMethodDescriptorByServerInfo(info.FullMethod)
		ctx = SetDescriptor(ctx, md)
		return handler(ctx, req)
	}
}

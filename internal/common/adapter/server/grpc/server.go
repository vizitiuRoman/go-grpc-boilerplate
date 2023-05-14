package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strconv"

	gateway "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server interface {
	MethodDescriptors
	RegisterOnShutdown(func())
	Register(descriptor *ServerDescriptor) error
	Start() error
	Stop()
}

type grpcServer struct {
	MethodDescriptors

	ctx            context.Context
	cfg            *Config
	logger         log.Logger
	grpcServer     *grpc.Server
	gatewayServer  *http.Server
	mux            *gateway.ServeMux
	opts           []grpc.DialOption
	grpcListenAddr string
}

func NewGRPC(
	ctx context.Context,
	logger log.Logger,
	cfg *Config,
	serveMux *gateway.ServeMux,
	// muxOptions []gateway.ServeMuxOption,
	// middlewares []Middleware,
) Server {
	//sort.Sort(ByOrder(middlewares))
	options := make([]grpc.ServerOption, 0)
	methodDescriptors := NewMethodDescriptors()
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		MakeMethodDescriptorInterceptor(methodDescriptors),
	}

	//for _, middleware := range middlewares {
	//	if middleware.Interceptor != nil {
	//		unaryInterceptors = append(unaryInterceptors, middleware.Interceptor)
	//	}
	//}

	options = append(options, grpc.MaxRecvMsgSize(cfg.MaxReceiveMessageLength), grpc.MaxSendMsgSize(cfg.MaxSendMessageLength))

	if len(unaryInterceptors) > 0 {
		options = append(options, grpc.ChainUnaryInterceptor(unaryInterceptors...))
	}
	//
	//for _, option := range muxOptions {
	//	option(serveMux)
	//}

	srv := &grpcServer{
		ctx:        ctx,
		cfg:        cfg,
		logger:     logger,
		grpcServer: grpc.NewServer(options...),
		mux:        serveMux,
		opts: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(cfg.MaxReceiveMessageLength),
				grpc.MaxCallSendMsgSize(cfg.MaxSendMessageLength),
			),
		},
		grpcListenAddr: net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.GrpcPort)),
		gatewayServer: &http.Server{
			Addr:    net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.HttpPort)),
			Handler: serveMux,
		},
		MethodDescriptors: methodDescriptors,
	}

	return srv
}

func (s *grpcServer) RegisterOnShutdown(f func()) {
	if s.gatewayServer != nil {
		s.gatewayServer.RegisterOnShutdown(f)
	}
}

func (s *grpcServer) Register(descriptor *ServerDescriptor) error {
	reflect.ValueOf(descriptor.GRPCRegistrar).Call([]reflect.Value{
		reflect.ValueOf(s.grpcServer),
		reflect.ValueOf(descriptor.Server),
	})

	s.MethodDescriptors.JoinServerDescriptor(descriptor)
	return descriptor.GRPCGatewayRegistrar(s.ctx, s.mux, fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.GrpcPort), s.opts)
}

func (s *grpcServer) startGRPC() error {
	listener, err := net.Listen("tcp", s.grpcListenAddr)
	if err != nil {
		return err
	}

	s.logger.Info("start grpc server", zap.String("addr", s.grpcListenAddr))

	return s.grpcServer.Serve(listener)
}

func (s *grpcServer) startGateway() error {
	s.logger.Info("start gateway", zap.String("addr", s.gatewayServer.Addr))

	return s.gatewayServer.ListenAndServe()
}

func (s *grpcServer) Start() error {
	errs := make(chan error, 1)
	go func() {
		errs <- s.startGRPC()
	}()
	go func() {
		errs <- s.startGateway()
	}()

	return <-errs
}

func (s *grpcServer) Stop() {
	defer s.logger.Debug("grpc server and gateway stopped")

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	if s.gatewayServer != nil {
		err := s.gatewayServer.Close()
		if err != nil {
			s.logger.Error("failed to close gateway server", zap.Error(err))
		}
	}
}

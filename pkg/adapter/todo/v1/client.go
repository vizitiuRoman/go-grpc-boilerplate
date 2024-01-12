package v1

import (
	"context"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/log"
	todov1 "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/gen/todo/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewTodoSVCClient(
	ctx context.Context,
	logger log.Logger,
	cfg *Config,
	dialOptions []grpc.DialOption,
	unaryInterceptors []grpc.UnaryClientInterceptor,
) (todov1.TodoSVCClient, error) {
	defaultDialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(c context.Context, s string) (net.Conn, error) {
			return net.DialTimeout("tcp", s, 10*time.Second)
		}),
	}

	if len(dialOptions) > 0 {
		defaultDialOptions = append(defaultDialOptions, dialOptions...)
	}

	if len(unaryInterceptors) > 0 {
		defaultDialOptions = append(defaultDialOptions, grpc.WithChainUnaryInterceptor(unaryInterceptors...))
	}

	f := &todoSVCFactory{
		ctx:    ctx,
		logger: logger.Named("todo-svc-client"),
		opts:   defaultDialOptions,
		cfg:    cfg,
		conns:  make(map[string]*grpc.ClientConn),
	}

	return f.CreateTodoSVCClient()
}

type todoSVCFactory struct {
	ctx context.Context
	mu  sync.Mutex

	logger log.Logger
	opts   []grpc.DialOption
	cfg    *Config
	conns  map[string]*grpc.ClientConn
}

func (f *todoSVCFactory) CreateTodoSVCClient() (todov1.TodoSVCClient, error) {
	conn, err := f.createConnection(net.JoinHostPort(f.cfg.Host, strconv.Itoa(f.cfg.GrpcPort)))
	if err != nil {
		return nil, err
	}

	return todov1.NewTodoSVCClient(conn), nil
}

func (f *todoSVCFactory) createConnection(target string) (*grpc.ClientConn, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.conns[target] != nil {
		return f.conns[target], nil
	}

	conn, err := grpc.Dial(target, f.opts...)
	if err != nil {
		f.logger.WithMethod(f.ctx, "createConnection").Error("unable to create client connection", zap.Error(err))
		return nil, err
	}
	f.conns[target] = conn

	return conn, nil
}

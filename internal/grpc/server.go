package grpcSrv

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/Cheasezz/balanceSrvc/internal/service"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type Config struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout"`
}

const (
	envLocal = "local"
)

type ServerAPI struct {
	blnc.UnimplementedBalanceServer
	Srvc *service.Service
}

type App struct {
	log    logger.Logger
	Server *grpc.Server
	port   int
}

func New(l logger.Logger, cfg Config, s *service.Service, env string) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			// logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			l.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(l), loggingOpts...),
	))

	blnc.RegisterBalanceServer(gRPCServer, &ServerAPI{Srvc: s})
	if env == envLocal {
		reflection.Register(gRPCServer)
	}

	return &App{
		log:    l,
		Server: gRPCServer,
		port:   cfg.Port,
	}
}

func InterceptorLogger(l logger.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, int(lvl), msg, fields...)
	})
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With("op", op)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", "addr", l.Addr().String())

	if err = a.Server.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) RunBufConn() func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(1024 * 1024)

	go func() {
		if err := a.Server.Serve(lis); err != nil {
			panic(err)
		}
	}()

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	return bufDialer
}

func (a *App) Close() {
	a.Server.GracefulStop()
}

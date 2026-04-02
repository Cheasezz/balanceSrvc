package grpcapp

import (
	"context"
	"fmt"
	"net"

	"github.com/Cheasezz/balanceSrvc/internal/config"
	grpcHndlrs "github.com/Cheasezz/balanceSrvc/internal/grpc"
	"github.com/Cheasezz/balanceSrvc/internal/service"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type App struct {
	log    logger.Logger
	Server *grpc.Server
	port   int
}

func New(l logger.Logger, cfg *config.Config, s *service.Service) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			// logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
	}
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(InterceptorLogger(l), loggingOpts...),
	))

	grpcHndlrs.Register(gRPCServer, l, s, cfg.Env)

	return &App{
		log:    l,
		Server: gRPCServer,
		port:   cfg.GRPC.Port,
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

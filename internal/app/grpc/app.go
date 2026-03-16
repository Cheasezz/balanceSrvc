package grpcapp

import (
	"fmt"
	"net"

	grpcHndlrs "github.com/Cheasezz/balanceSrvc/internal/grpc"
	"github.com/Cheasezz/balanceSrvc/internal/service"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	"google.golang.org/grpc"
)

type App struct {
	log        logger.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(l logger.Logger, p int, s *service.Service) *App {
	gRPCServer := grpc.NewServer()

	grpcHndlrs.Register(gRPCServer, l, s)

	return &App{
		log:        l,
		gRPCServer: gRPCServer,
		port:       p,
	}
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

	if err = a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
}

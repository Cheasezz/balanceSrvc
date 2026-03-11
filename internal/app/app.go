package app

import (
	"log/slog"

	grpcapp "github.com/Cheasezz/balanceSrvc/internal/app/grpc"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(l *slog.Logger, p int, pbUrl string) *App {
	grpcApp := grpcapp.New(l, p)

	return &App{GRPCSrv: grpcApp}
}

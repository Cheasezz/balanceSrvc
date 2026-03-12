package app

import (
	grpcapp "github.com/Cheasezz/balanceSrvc/internal/app/grpc"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(l logger.Logger, p int, pbUrl string) *App {
	grpcApp := grpcapp.New(l, p)

	return &App{GRPCSrv: grpcApp}
}

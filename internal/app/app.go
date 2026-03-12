package app

import (
	grpcapp "github.com/Cheasezz/balanceSrvc/internal/app/grpc"
	"github.com/Cheasezz/balanceSrvc/internal/repo"
	"github.com/Cheasezz/balanceSrvc/internal/service"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(l logger.Logger, p int, pbUrl string) *App {
	db := repo.New()

	srvc := service.New(l, db)

	grpcApp := grpcapp.New(l, p, srvc)

	return &App{GRPCSrv: grpcApp}
}

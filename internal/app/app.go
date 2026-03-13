package app

import (
	grpcapp "github.com/Cheasezz/balanceSrvc/internal/app/grpc"
	"github.com/Cheasezz/balanceSrvc/internal/repo"
	"github.com/Cheasezz/balanceSrvc/internal/service"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
)

type App struct {
	GRPCSrv *grpcapp.App
	DB      *pgx5.Pgx
	l       logger.Logger
}

func New(l logger.Logger, p int, pbUrl string) *App {
	const op = "app.New"

	db, err := pgx5.New(pbUrl)
	if err != nil {
		l.With("op", op).Error(err.Error())
		panic("pgx5 can't create")
	}

	repo := repo.New(db)

	srvc := service.New(l, repo)

	grpcApp := grpcapp.New(l, p, srvc)

	return &App{GRPCSrv: grpcApp, DB: db, l: l}
}

func (a *App) Close() {
	const op = "app.Close"
	log := a.l.With("op", op)

	log.Info("stopping gRPC server")
	a.GRPCSrv.Stop()

	log.Info("stopping DB")
	a.DB.Close()
}

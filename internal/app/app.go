package app

import (
	"context"

	trxtyperegistry "github.com/Cheasezz/balanceSrvc/internal/app/trxTypeRegistry"
	"github.com/Cheasezz/balanceSrvc/internal/config"
	grpcSrv "github.com/Cheasezz/balanceSrvc/internal/grpc"
	"github.com/Cheasezz/balanceSrvc/internal/repo"
	"github.com/Cheasezz/balanceSrvc/internal/service"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
)

type App struct {
	GRPCApp *grpcSrv.App
	db      *pgx5.Pgx
	l       logger.Logger
}

func New(l logger.Logger, cfg *config.Config) *App {
	const op = "app.New"
	log := l.With("op", op)

	db, err := pgx5.New(cfg.PG)
	if err != nil {
		log.Error(err.Error())
		panic("pgx5 can't create")
	}

	repo := repo.New(db)

	dbTrxTypes, err := repo.Trx.GetAllTypesInfo(context.Background())
	if err != nil {
		log.Error(err.Error())
		panic("can't collect db transaction types")
	}

	registry, err := trxtyperegistry.New(dbTrxTypes)
	if err != nil {
		log.Error(err.Error())
		panic("can't create transaction types registry")
	}

	srvc := service.New(l, repo, registry)

	grpcApp := grpcSrv.New(l, cfg.GRPC, srvc, cfg.Env)

	return &App{GRPCApp: grpcApp, db: db, l: l}
}

func (a *App) Close() {
	const op = "app.Close"
	log := a.l.With("op", op)

	log.Info("stopping gRPC server")
	a.GRPCApp.Close()

	log.Info("stopping DB")
	a.db.Close()
}

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Cheasezz/balanceSrvc/internal/app"
	"github.com/Cheasezz/balanceSrvc/internal/config"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	log.Info("starting application")

	application := app.New(log, cfg.GRPC.Port, cfg.PG.URL)

	go application.GRPCSrv.MustRun()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()
	log.Info("Gracefully stopped")
}

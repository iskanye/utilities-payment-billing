package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/iskanye/utilities-payment-billing/internal/app"
	"github.com/iskanye/utilities-payment-billing/internal/config"
	pkgConfig "github.com/iskanye/utilities-payment-utils/pkg/config"
	"github.com/iskanye/utilities-payment-utils/pkg/logger"
)

func main() {
	cfg := pkgConfig.MustLoad[config.Config]()
	cfg.LoadEnv()

	log := logger.SetupPrettySlog()
	app := app.New(
		log,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Term,
		cfg.GRPC.Port,
	)

	go func() {
		app.GRPCServer.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GRPCServer.Stop()
	log.Info("Gracefully stopped")
}

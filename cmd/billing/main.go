package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/iskanye/utilities-payment-billing/internal/app"
	"github.com/iskanye/utilities-payment-billing/internal/config"
	pkgConfig "github.com/iskanye/utilities-payment/pkg/config"
	"github.com/iskanye/utilities-payment/pkg/logger"
)

func main() {
	cfg := pkgConfig.MustLoad[config.Config](pkgConfig.NoModyfing)
	log := logger.SetupPrettySlog()
	app := app.New(
		log,
		cfg.Postgre.User,
		cfg.Postgre.Password,
		cfg.Postgre.DBName,
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

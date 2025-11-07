package app

import (
	"log/slog"

	"github.com/iskanye/utilities-payment-billing/internal/app/grpc"
	"github.com/iskanye/utilities-payment-billing/internal/service/billing"
	"github.com/iskanye/utilities-payment-billing/internal/storage"
)

type App struct {
	GRPCServer *grpc.App
}

func New(
	log *slog.Logger,
	dbHost string,
	dbPort int,
	dbUser string,
	dbPassword string,
	dbName string,
	term int,
	grpcPort int,
) *App {
	storage, err := storage.New(dbHost, dbPort, dbUser, dbPassword, dbName, term)
	if err != nil {
		panic(err)
	}

	billingService := billing.New(log, storage, storage, storage)
	grpcApp := grpc.New(log, billingService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}

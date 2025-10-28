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
	user string,
	password string,
	dbName string,
	term int,
	grpcPort int,
) *App {
	storage, err := storage.New(user, password, dbName, term)
	if err != nil {
		panic(err)
	}

	billingService := billing.New(log, storage, storage)
	grpcApp := grpc.New(log, billingService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}

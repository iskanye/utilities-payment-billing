package app

import (
	"log/slog"

	"github.com/iskanye/utilities-payment-billing/internal/app/grpc"
)

type App struct {
	GRPCServer *grpc.App
}

func New(
	log *slog.Logger,
	grpcPort int,
) *App {
	return &App{}
}

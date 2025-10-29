package suite

import (
	"context"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/iskanye/utilities-payment-billing/internal/config"
	"github.com/iskanye/utilities-payment-proto/billing"
	pkgConfig "github.com/iskanye/utilities-payment/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	Cfg           *config.Config        // Конфигурация приложения
	BillingClient billing.BillingClient // Клиент для взаимодействия с gRPC-сервером
}

const (
	grpcHost = "localhost"
)

// New creates new test suite.
//
// TODO: for pipeline tests we need to wait for app is ready
func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := pkgConfig.MustLoadPath(configPath(), func(t *config.Config) {})

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		Cfg:           cfg,
		BillingClient: billing.NewBillingClient(cc),
	}
}

func configPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "../config/tests.yaml"
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}

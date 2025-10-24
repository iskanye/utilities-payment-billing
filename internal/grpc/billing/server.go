package billing

import (
	"context"

	"google.golang.org/grpc"

	protoAuth "github.com/iskanye/utilities-payment-proto/billing"
)

type serverAPI struct {
	protoAuth.UnimplementedBillingServer
	auth Billing
}

type Billing interface {
	AddBill(
		ctx context.Context,
		address string,
		amount int,
	) (bill_id int64, err error)
}

func Register(gRPCServer *grpc.Server, billing Billing) {
	protoAuth.RegisterBillingServer(gRPCServer, &serverAPI{auth: billing})
}

func (s *serverAPI) AddBill(
	ctx context.Context,
	in *protoAuth.Bill,
) (*protoAuth.BillResponse, error) {
	return &protoAuth.BillResponse{}, nil
}

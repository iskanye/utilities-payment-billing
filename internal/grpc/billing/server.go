package billing

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	protoBilling "github.com/iskanye/utilities-payment-proto/billing"
)

type serverAPI struct {
	protoBilling.UnimplementedBillingServer
	billing Billing
}

type Billing interface {
	AddBill(
		ctx context.Context,
		address string,
		amount int,
	) (bill_id int64, err error)
}

func Register(gRPCServer *grpc.Server, billing Billing) {
	protoBilling.RegisterBillingServer(gRPCServer, &serverAPI{billing: billing})
}

func (s *serverAPI) AddBill(
	ctx context.Context,
	in *protoBilling.Bill,
) (*protoBilling.BillResponse, error) {
	if in.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address is required")
	}
	if in.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be positive")
	}

	billId, err := s.billing.AddBill(ctx, in.GetAddress(), int(in.GetAmount()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &protoBilling.BillResponse{
		BillId: billId,
	}, nil
}

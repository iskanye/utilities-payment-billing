package billing

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	protoBilling "github.com/iskanye/utilities-payment-proto/billing"
	"github.com/iskanye/utilities-payment/pkg/models"
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
	) (int64, error)
	GetBills(
		ctx context.Context,
		address string,
	) ([]models.Bill, error)
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

func (s *serverAPI) GetBills(
	in *protoBilling.BillsRequest,
	stream grpc.ServerStreamingServer[protoBilling.Bill],
) error {
	if in.Address == "" {
		return status.Error(codes.InvalidArgument, "address is required")
	}

	bills, err := s.billing.GetBills(stream.Context(), in.GetAddress())
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, i := range bills {
		if err = stream.SendMsg(i); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}

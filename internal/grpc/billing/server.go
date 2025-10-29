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
		userID int64,
	) (int64, error)
	GetBills(
		ctx context.Context,
		userID int64,
	) ([]models.Bill, error)
	PayBill(
		ctx context.Context,
		billId int64,
	) error
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
	if in.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	billId, err := s.billing.AddBill(ctx, in.GetAddress(), int(in.GetAmount()), in.GetUserId())
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
	if in.UserId == 0 {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	bills, err := s.billing.GetBills(stream.Context(), in.GetUserId())
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, i := range bills {
		id := i.ID
		dueDate := i.DueDate
		bill := &protoBilling.Bill{
			BillId:  &id,
			Address: i.Address,
			Amount:  int32(i.Amount),
			UserId:  i.UserID,
			DueDate: &dueDate,
		}
		if err = stream.SendMsg(bill); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}

func (s *serverAPI) PayBill(
	ctx context.Context,
	in *protoBilling.PayRequest,
) (*protoBilling.PayResponse, error) {
	if in.BillId == 0 {
		return nil, status.Error(codes.InvalidArgument, "bill_id is required")
	}

	err := s.billing.PayBill(ctx, in.GetBillId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &protoBilling.PayResponse{}, nil
}

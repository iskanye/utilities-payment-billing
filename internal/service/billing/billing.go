package billing

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/iskanye/utilities-payment-billing/internal/storage"
	"github.com/iskanye/utilities-payment/pkg/logger"
	"github.com/iskanye/utilities-payment/pkg/models"
)

type Billing struct {
	log           *slog.Logger
	billCreator   BillCreator
	billsProvider BillsProvider
	billPayer     BillPayer
}

type BillCreator interface {
	CreateBill(
		ctx context.Context,
		address string,
		amount int,
		userID int64,
	) (int64, error)
}

type BillsProvider interface {
	GetBills(
		ctx context.Context,
		userID int64,
	) ([]models.Bill, error)
}

type BillPayer interface {
	PayBill(
		ctx context.Context,
		billId int64,
	) error
}

func New(
	log *slog.Logger,
	billCreator BillCreator,
	billsProvider BillsProvider,
	billPayer BillPayer,
) *Billing {
	return &Billing{
		log:           log,
		billCreator:   billCreator,
		billsProvider: billsProvider,
		billPayer:     billPayer,
	}
}

func (b *Billing) AddBill(
	ctx context.Context,
	address string,
	amount int,
	userID int64,
) (int64, error) {
	const op = "Billing.AddBill"

	log := b.log.With(
		slog.String("op", op),
		slog.String("address", address),
	)

	log.Info("attempting to create bill")

	billId, err := b.billCreator.CreateBill(ctx, address, amount, userID)
	if err != nil {
		log.Error("failed to create bill", logger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("created bill successfully",
		slog.Int64("bill_id", billId),
	)

	return billId, nil
}

func (b *Billing) GetBills(
	ctx context.Context,
	userID int64,
) ([]models.Bill, error) {
	const op = "Billing.GetBills"

	log := b.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("attempting to get bill")

	bills, err := b.billsProvider.GetBills(ctx, userID)
	if err != nil {
		log.Error("failed to get bill", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully got bills")

	return bills, nil
}

func (b *Billing) PayBill(
	ctx context.Context,
	billId int64,
) error {
	const op = "Billing.PayBill"

	log := b.log.With(
		slog.String("op", op),
		slog.Int64("bill_id", billId),
	)

	log.Info("attempting to pay the bill")

	err := b.billPayer.PayBill(ctx, billId)
	if err != nil {
		if errors.Is(err, storage.ErrBillNotFound) {
			log.Warn("bill not found", logger.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to get bill", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully payed the bill")

	return nil
}

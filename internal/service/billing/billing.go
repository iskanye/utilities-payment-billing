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
}

type BillCreator interface {
	CreateBill(
		ctx context.Context,
		address string,
		amount int,
	) (int64, error)
}

type BillsProvider interface {
	GetBills(
		ctx context.Context,
		address string,
	) ([]models.Bill, error)
}

func New(
	log *slog.Logger,
	billCreator BillCreator,
	billsProvider BillsProvider,
) *Billing {
	return &Billing{
		log:           log,
		billCreator:   billCreator,
		billsProvider: billsProvider,
	}
}

func (b *Billing) AddBill(
	ctx context.Context,
	address string,
	amount int,
) (int64, error) {
	const op = "Billing.AddBill"

	log := b.log.With(
		slog.String("op", op),
		slog.String("address", address),
	)

	log.Info("attempting to create bill")

	billId, err := b.billCreator.CreateBill(ctx, address, amount)
	if err != nil {
		log.Error("failed to create bill: ", logger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return billId, nil
}

func (b *Billing) GetBills(
	ctx context.Context,
	address string,
) ([]models.Bill, error) {
	const op = "Billing.GetBills"

	log := b.log.With(
		slog.String("op", op),
		slog.String("address", address),
	)

	log.Info("attempting to get bill")

	bills, err := b.billsProvider.GetBills(ctx, address)
	if err != nil {
		if errors.Is(err, storage.ErrBillsNotFound) {
			log.Warn("bills not found: ", logger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to get bill: ", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return bills, nil
}

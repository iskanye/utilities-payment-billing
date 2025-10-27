package billing

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iskanye/utilities-payment/pkg/logger"
	"github.com/iskanye/utilities-payment/pkg/models"
)

type Billing struct {
	log         *slog.Logger
	billCreator BillCreator
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
) *Billing {
	return &Billing{
		log:         log,
		billCreator: billCreator,
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
	const op = "Billing.GetBill"

	log := b.log.With(
		slog.String("op", op),
		slog.String("address", address),
	)

	log.Info("attempting to get bill")

	// TO BE IMPLEMENTED
	/*bill, err := b.billProvider.GetBill(ctx, address)
	if err != nil {
		if errors.Is(err, storage.ErrBillNotFound) {
			log.Warn("bill not found: ", logger.Err(err))
			return models.Bill{}, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to get bill: ", logger.Err(err))
		return models.Bill{}, fmt.Errorf("%s: %w", op, err)
	}*/

	return nil, nil
}

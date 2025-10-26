package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/iskanye/utilities-payment-billing/internal/config"
	"github.com/iskanye/utilities-payment/pkg/models"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(
	c *config.Config,
) (*Storage, error) {
	const op = "storage.postgre.New"

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		c.Postgre.User,
		c.Postgre.Password,
		c.Postgre.DBName,
	)
	db, err := sql.Open("postgre", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) CreateBill(
	ctx context.Context,
	address string,
	amount int,
) (int64, error) {
	const op = "storage.postgre.CreateBill"

	stmt, err := s.db.Prepare("INSERT INTO bills(address, amount, due_date) VALUES(?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, address, amount, time.Now().AddDate(0, 1, 0))
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetBill(
	ctx context.Context,
	billId int64,
) (models.Bill, error) {
	const op = "storage.postgre.CreateBill"

	stmt, err := s.db.Prepare("SELECT address, amount, due_date FROM bills WHERE id = ?")
	if err != nil {
		return models.Bill{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, billId)

	var bill models.Bill
	err = row.Scan(&bill.ID, &bill.Address, &bill.DueDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Bill{}, fmt.Errorf("%s: %w", op, ErrBillNotFound)
		}

		return models.Bill{}, fmt.Errorf("%s: %w", op, err)
	}

	return bill, nil
}

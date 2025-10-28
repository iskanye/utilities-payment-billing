package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/iskanye/utilities-payment/pkg/models"
	_ "github.com/lib/pq"
)

type Storage struct {
	db   *sql.DB
	term int // in Months
}

func New(
	user string,
	password string,
	dbName string,
	term int,
) (*Storage, error) {
	const op = "storage.postgres.New"

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		user, password, dbName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db:   db,
		term: term,
	}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) CreateBill(
	ctx context.Context,
	address string,
	amount int,
) (int64, error) {
	const op = "storage.postgres.CreateBill"

	stmt, err := s.db.Prepare("INSERT INTO bills(address, amount, due_date) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res := stmt.QueryRowContext(ctx, address, amount, time.Now().AddDate(0, s.term, 0))

	var id int64
	err = res.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetBills(
	ctx context.Context,
	address string,
) ([]models.Bill, error) {
	const op = "storage.postgres.GetBills"

	stmt, err := s.db.Prepare("SELECT id, address, amount, due_date FROM bills WHERE address = $1;")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()
	bills := make([]models.Bill, 0)

	for rows.Next() {
		var bill models.Bill
		err = rows.Scan(&bill.ID, &bill.Address, &bill.Amount, &bill.DueDate)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%s: %w", op, ErrBillsNotFound)
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		bills = append(bills, bill)
	}

	return bills, nil
}

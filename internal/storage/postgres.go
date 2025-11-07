package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/iskanye/utilities-payment-utils/pkg/models"
	_ "github.com/lib/pq"
)

type Storage struct {
	db   *sql.DB
	term int // in Months
}

func New(
	host string,
	port int,
	user string,
	password string,
	dbName string,
	term int,
) (*Storage, error) {
	const op = "storage.postgres.New"

	connStr := fmt.Sprintf(
		"host=%s post=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
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
	userID int64,
) (int64, error) {
	const op = "storage.postgres.CreateBill"

	stmt, err := s.db.Prepare("INSERT INTO bills(address, amount, user_id, due_date) VALUES($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res := stmt.QueryRowContext(ctx, address, amount, userID, time.Now().AddDate(0, s.term, 0))

	var id int64
	err = res.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetBills(
	ctx context.Context,
	userId int64,
) ([]models.Bill, error) {
	const op = "storage.postgres.GetBills"

	stmt, err := s.db.Prepare("SELECT id, address, amount, user_id, due_date FROM bills WHERE user_id = $1;")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()
	var bills []models.Bill

	for rows.Next() {
		var bill models.Bill
		err = rows.Scan(&bill.ID, &bill.Address, &bill.Amount, &bill.UserID, &bill.DueDate)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		bills = append(bills, bill)
	}

	return bills, nil
}

func (s *Storage) GetBill(
	ctx context.Context,
	billID int64,
) (models.Bill, error) {
	const op = "storage.postgres.GetBill"

	stmt, err := s.db.Prepare("SELECT id, address, amount, user_id, due_date FROM bills WHERE id = $1;")
	if err != nil {
		return models.Bill{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, billID)

	var bill models.Bill
	err = row.Scan(&bill.ID, &bill.Address, &bill.Amount, &bill.UserID, &bill.DueDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Bill{}, fmt.Errorf("%s: %w", op, ErrBillNotFound)
		}

		return models.Bill{}, fmt.Errorf("%s: %w", op, err)
	}

	return bill, nil
}

func (s *Storage) PayBill(
	ctx context.Context,
	billId int64,
) error {
	const op = "storage.postgres.PayBill"

	stmt, err := s.db.Prepare("DELETE FROM bills WHERE id = $1;")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, billId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rows == 0 {
		return fmt.Errorf("%s: %w", op, ErrBillNotFound)
	}

	return nil
}

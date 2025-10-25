package storage

import (
	"database/sql"
	"fmt"

	"github.com/iskanye/utilities-payment-billing/internal/config"
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

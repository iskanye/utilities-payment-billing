package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/iskanye/utilities-payment-billing/internal/config"
	pkgConfig "github.com/iskanye/utilities-payment-utils/pkg/config"
)

func main() {
	cfg := pkgConfig.MustLoad[config.Config]()

	uri := fmt.Sprintf("%s:%s@%s:%d/%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
	)

	var migrationsPath, migrationsTable string
	var clear bool

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.BoolVar(&clear, "clear", false, "use down migrations")
	flag.Parse()

	if uri == "" {
		panic("uri is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("postgres://%s?x-migrations-table=%s&sslmode=disable", uri, migrationsTable),
	)
	if err != nil {
		panic(err)
	}
	if clear {
		mustMigrate(m.Down())
	} else {
		mustMigrate(m.Up())
	}

	fmt.Println("migrations applied")
}

func mustMigrate(err error) {
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}

		panic(err)
	}
}

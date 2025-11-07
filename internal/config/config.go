package config

import "time"

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	GRPC     GRPCConfig     `yaml:"grpc"`
	Term     int            `yaml:"term"` // in Months
}

type PostgresConfig struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

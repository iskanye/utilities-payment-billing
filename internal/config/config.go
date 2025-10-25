package config

import "time"

type Config struct {
	Postgre PostgreConfig `yaml:"postgre"`
	GRPC    GRPCConfig    `yaml:"grpc"`
}

type PostgreConfig struct {
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

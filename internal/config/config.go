package config

import "os"

type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

func FromEnv() Config {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://app:app@localhost:5432/app?sslmode=disable"
	}
	return Config{
		HTTPAddr:    addr,
		DatabaseURL: dsn,
	}
}

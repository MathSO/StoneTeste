package database

import (
	"context"
	"fmt"
	"server/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool

func connect() (*pgxpool.Pool, error) {
	if conn != nil {
		return conn, nil
	}

	connectStr := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`, config.POSTGRES_DB_HOST, config.POSTGRES_DB_PORT, config.POSTGRES_DB_USER, config.POSTGRES_DB_PASS, config.POSTGRES_DB_NAME)
	cfg, err := pgxpool.ParseConfig(connectStr)
	if err != nil {
		return nil, err
	}

	conn, err = pgxpool.NewWithConfig(context.Background(), cfg)
	return conn, err
}

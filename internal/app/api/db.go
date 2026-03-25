package api

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func openDB(cfg Config) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConns = int32(cfg.DB.MaxOpenConns)
	dbConfig.MinConns = int32(cfg.DB.MaxIdleConns)
	dbConfig.MaxConnIdleTime = cfg.DB.MaxIdleTime
	conn, err := pgxpool.New(context.Background(), cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

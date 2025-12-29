package db

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PgxIface interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

func Connect(log *zap.Logger) (*pgxpool.Pool, error) {

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Error("db URL is empty")
		return nil, errors.New("db URL is empty")
	}

	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Error("cant parse db config", zap.Error(err))
		return nil, err
	}

	// db pool config
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Error("cant create db pool")
		return nil, err
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Error("cant ping to db")
		return nil, err
	}

	log.Info("connected to DB", zap.Time("At", time.Now()))
	return dbPool, nil
}

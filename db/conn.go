package db

import (
	"context"
	"fmt"
	"time"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
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

func Connect(log *zap.Logger, config utils.DatabaseConfiguration) (*pgxpool.Pool, error) {
	// "user password dbname sslmode host"
	dbUrl := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s",
		config.UserName, config.Password, config.DBName, config.SSLMode, config.HostName)
	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Error("cant parse db config", zap.Error(err))
		return nil, err
	}

	// db pool config
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Error("cant create db pool", zap.Error(err))
		return nil, err
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Error("cant ping to db", zap.Error(err))
		return nil, err
	}

	log.Info("connected to DB", zap.Time("At", time.Now()))
	return dbPool, nil
}

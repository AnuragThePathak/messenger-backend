package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func NewPgxPool(ctx context.Context, connString string, logger pgx.Logger,
	logLevel pgx.LogLevel) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	conf.ConnConfig.Logger = logger

	if logLevel != 0 {
		conf.ConnConfig.LogLevel = logLevel
	}

	conf.LazyConnect = true

	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("pgx connection error: %w", err)
	}
	return pool, err
}

func LogLevelFromEnv() (pgx.LogLevel, error) {
	if level := os.Getenv("PGX_LOG_LEVEL"); level != "" {
		l, err := pgx.LogLevelFromString(level)
		if err != nil {
			return pgx.LogLevelDebug, fmt.Errorf("pgx configuration: %w", err)
		}
		return l, nil
	}
	return pgx.LogLevelInfo, nil
}

func GetLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}

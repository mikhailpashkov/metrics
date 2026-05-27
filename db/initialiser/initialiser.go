package initialiser

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	dbmetrics "github.com/mikhailpashkov/metrics/db/metrics"
	"github.com/mikhailpashkov/metrics/db/migrations"
)

type Queries struct {
	Metrics *dbmetrics.Queries
}

func InitialiseDB(logger *slog.Logger, databaseDSN string) (*Queries, error) {
	logger.Debug("connect to db")

	pgxPool, err := pgxpool.New(context.Background(), databaseDSN)
	if err != nil {
		logger.Error("failed to connect to DB", "err", err.Error())
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}
	defer pgxPool.Close()

	logger.Debug("test db connection")

	dbPingTimeoutCtx, cancelFunc := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancelFunc()

	if err = pgxPool.Ping(dbPingTimeoutCtx); err != nil {
		logger.Error("failed to ping DB", "err", err.Error())
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	err = migrations.RunMigrations(logger, pgxPool)
	if err != nil {
		logger.Error("failed to run migrations", "err", err.Error())
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &Queries{
		Metrics: dbmetrics.New(pgxPool),
	}, nil
}

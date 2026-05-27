package repository

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var retryDelays = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	5 * time.Second,
}

func isPGRetriable(err error) bool {
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return false
	}

	switch pgErr.Code {
	case
		pgerrcode.ConnectionException,
		pgerrcode.ConnectionFailure,
		pgerrcode.ConnectionDoesNotExist,
		pgerrcode.SerializationFailure,
		pgerrcode.DeadlockDetected:
		return true
	}

	return false
}

func PGRetry(ctx context.Context, fn func() error, logger *slog.Logger) error {
	var err error

	err = fn()
	if err == nil || !isPGRetriable(err) {
		return err
	}

	for i, delay := range retryDelays {
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			logger.Debug("ctx.Done()", "err", ctx.Err())
			return ctx.Err()
		}

		err = fn()
		if err == nil {
			return nil
		}

		if !isPGRetriable(err) {
			logger.Debug("not pg retriable error", "err", err)
			return err
		}

		logger.Debug("will retry after delay", "retry", i+1, "delay", delay)
	}

	return err
}

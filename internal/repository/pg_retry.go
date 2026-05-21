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

func IsPGRetriable(err error) bool {
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

	for i, delay := range retryDelays {
		err = fn()
		if err == nil {
			return nil
		}

		if !IsPGRetriable(err) {
			logger.Debug("not pg retriable error", "err", err)
			return err
		}

		logger.Debug("will retry after delay", "retry", i+1, "delay", delay)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			logger.Debug("ctx.Done()", "err", ctx.Err())
			return ctx.Err()
		}
	}

	return err
}

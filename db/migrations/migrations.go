package migrations

import (
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var FS embed.FS

func RunMigrations(logger *slog.Logger, pgxPool *pgxpool.Pool) error {
	logger.Debug("run migrations")

	goose.SetLogger(newGooseLogger(logger))
	goose.SetBaseFS(FS)

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("failed to goose.SetDialect postgres", "err", err.Error())
		return err
	}

	sqlDB := stdlib.OpenDBFromPool(pgxPool)
	if err := goose.Up(sqlDB, "."); err != nil {
		logger.Error("failed to goose.Up migrations", "err", err.Error())
		return err
	}
	return nil
}

type gooseLogger struct {
	*slog.Logger
}

func newGooseLogger(logger *slog.Logger) *gooseLogger {
	return &gooseLogger{logger}
}

func (gl *gooseLogger) Fatalf(format string, v ...any) {
	gl.Logger.Error(fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (gl *gooseLogger) Printf(format string, v ...any) {
	gl.Logger.Info(fmt.Sprintf(format, v...))
}

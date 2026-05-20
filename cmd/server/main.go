package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/mikhailpashkov/metrics/internal/handler"
	"github.com/mikhailpashkov/metrics/internal/handler/middleware"
	"github.com/mikhailpashkov/metrics/internal/repository"
	"github.com/mikhailpashkov/metrics/internal/service"
	"github.com/mikhailpashkov/metrics/internal/utils"
)

const (
	LoggerNameKey = "slog_logger"
)

func main() {
	// Logger /////////////////////////
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	logger.Info("SERVER")

	// Params /////////////////////////
	var addr string
	var storeInterval int
	var fileStoragePath string
	var restore bool
	var databaseDSN string

	utils.GetParams([]utils.Param{
		&utils.StringParam{
			EnvName:       "ADDRESS",
			FlagName:      "a",
			FlagUsage:     "Server address",
			Default:       "localhost:8080",
			ValueConsumer: func(v string) { addr = v },
		},
		&utils.IntParam{
			EnvName:       "STORE_INTERVAL",
			FlagName:      "i",
			FlagUsage:     "Store interval",
			Default:       300,
			ValueConsumer: func(v int) { storeInterval = v },
		},
		&utils.StringParam{
			EnvName:       "FILE_STORAGE_PATH",
			FlagName:      "f",
			FlagUsage:     "File storage path",
			Default:       "./backup.json",
			ValueConsumer: func(v string) { fileStoragePath = v },
		},
		&utils.BoolParam{
			EnvName:       "RESTORE",
			FlagName:      "r",
			FlagUsage:     "Restore backup on startup",
			Default:       true,
			ValueConsumer: func(v bool) { restore = v },
		},
		&utils.StringParam{
			EnvName:       "DATABASE_DSN",
			FlagName:      "d",
			FlagUsage:     "Database connection string",
			Default:       "postgresql://username:password@localhost:5432/default_database",
			ValueConsumer: func(v string) { databaseDSN = v },
		},
	})

	logger.Debug("params read",
		"serverAddr", addr,
		"storeInterval", storeInterval,
		"fileStoragePath", fileStoragePath,
		"restore", restore,
		"databaseDSN", databaseDSN,
	)

	// Dependencies ///////////////////
	logger.Debug("init dependencies")
	metricsRepository := repository.NewMetricsMemoryRepository()
	backupRepository := repository.NewFileBackupRepository(fileStoragePath)
	eventService := service.NewInMemoryEventService(logger.With(LoggerNameKey, "service.NewInMemoryEventService"))
	metricsService := service.NewMetricsService(
		logger.With(LoggerNameKey, "service.MetricsService"),
		metricsRepository,
		eventService,
	)
	backupService := service.NewBackupService(
		logger.With(LoggerNameKey, "service.BackupService"),
		metricsService,
		backupRepository,
		eventService,
	)

	// Backup /////////////////////////
	if restore {
		logger.Debug("restore from backup")
		err := backupService.Restore(context.Background())
		if err != nil {
			logger.Error("failed to restore from backup", "err", err)
			os.Exit(1)
		}
	}
	logger.Debug("setup runtime backup")
	err := backupService.SetupBackup(context.Background(), storeInterval)
	if err != nil {
		logger.Error("failed to setup backup", "err", err)
		os.Exit(1)
	}

	// Server /////////////////////////
	logger.Debug("setup server")

	r := chi.NewRouter()

	// наверняка, хорошей идеей будет использовать github.com/go-chi/chi/v5/middleware,
	// но в учебных целях используем самодельные
	r.Use(middleware.WithLogging(logger.With(LoggerNameKey, "middleware.WithLogging")))
	r.Use(middleware.WithGZIPSupport(logger.With(LoggerNameKey, "middleware.WithGZIPSupport")))

	// для фикса автотестов в iter7: там, зачем-то, в конце слеши приделали на клиенте
	r.Use(chimiddleware.StripSlashes)

	r.Handle("/", handler.NewMetricsRootHandler(
		logger.With(LoggerNameKey, "handler.MetricsRootHandler"),
		metricsService,
	))

	r.Handle("/value", handler.NewGetMetricsHandler(
		logger.With(LoggerNameKey, "handler.GetMetricsHandler"),
		metricsService,
	))
	r.Handle("/value/{type}/{name}", handler.NewGetMetricsPathParamsHandler(
		logger.With(LoggerNameKey, "handler.GetMetricsPathParamsHandler"),
		metricsService,
	))

	r.Handle("/update", handler.NewUpdateMetricsHandler(
		logger.With(LoggerNameKey, "handler.UpdateMetricsHandler"),
		metricsService,
	))
	r.Handle("/update/{type}/{name}/{value}", handler.NewUpdateMetricsPathParamsHandler(
		logger.With(LoggerNameKey, "handler.UpdateMetricsPathParamsHandler"),
		metricsService,
	))

	err = http.ListenAndServe(addr, r)
	if err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

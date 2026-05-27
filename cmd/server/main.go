package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/mikhailpashkov/metrics/db/initialiser"
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
			Default:       false,
			ValueConsumer: func(v bool) { restore = v },
		},
		&utils.StringParam{
			EnvName:       "DATABASE_DSN",
			FlagName:      "d",
			FlagUsage:     "Database connection string",
			Default:       "",
			ValueConsumer: func(v string) { databaseDSN = v },
		},
	})

	logger.Debug("params read",
		"serverAddr", addr,
		"storeInterval", storeInterval,
		"fileStoragePath", fileStoragePath,
		"restore", restore,
		"len(databaseDSN)", len(databaseDSN), // dont log sensitive data
	)

	// Database ///////////////////////
	wantDB := len(databaseDSN) != 0
	var dbQueries initialiser.Queries
	var err error
	if wantDB {
		queries, err := initialiser.InitialiseDB(logger.With(LoggerNameKey, "initialiser.InitialiseDB"), databaseDSN)
		if err != nil {
			logger.Error("failed to initialise db queries", "err", err)
			os.Exit(1)
		}
		dbQueries = *queries
	} else {
		logger.Debug("empty databaseDSN, skip connect to db")
	}

	// Dependencies ///////////////////
	logger.Debug("init dependencies")

	var metricsRepository service.MetricsRepository
	if wantDB {
		metricsRepository = repository.NewMetricsDBRepository(dbQueries.Metrics, logger.With(LoggerNameKey, "repository.MetricsDBRepository"))
	} else {
		metricsRepository = repository.NewMetricsMemoryRepository()
	}

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
	err = backupService.SetupBackup(context.Background(), storeInterval)
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

	r.Get("/", handler.NewMetricsRootHandlerFunc(
		logger.With(LoggerNameKey, "handler.MetricsRootHandler"),
		metricsService,
	))

	r.Post("/value", handler.NewGetMetricsHandlerFunc(
		logger.With(LoggerNameKey, "handler.GetMetricsHandler"),
		metricsService,
	))
	r.Get("/value/{type}/{name}", handler.NewGetMetricsPathParamsHandlerFunc(
		logger.With(LoggerNameKey, "handler.GetMetricsPathParamsHandler"),
		metricsService,
	))

	r.Post("/update", handler.NewUpdateMetricsHandlerFunc(
		logger.With(LoggerNameKey, "handler.UpdateMetricsHandler"),
		metricsService,
	))
	r.Post("/update/{type}/{name}/{value}", handler.NewUpdateMetricsPathParamsHandlerFunc(
		logger.With(LoggerNameKey, "handler.UpdateMetricsPathParamsHandler"),
		metricsService,
	))
	r.Post("/updates", handler.NewUpdateMetricsBatchHandlerFunc(
		logger.With(LoggerNameKey, "handler.NewUpdateMetricsBatchHandler"),
		metricsService,
	))

	if wantDB {
		r.Get("/ping", handler.NewDBPingHandlerFunc(
			logger.With(LoggerNameKey, "handler.DBPingHandler"),
			dbQueries.Metrics,
		))
	}

	err = http.ListenAndServe(addr, r)
	if err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

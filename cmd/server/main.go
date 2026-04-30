package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/mikhailpashkov/metrics/internal/handler"
	"github.com/mikhailpashkov/metrics/internal/handler/middleware"
	"github.com/mikhailpashkov/metrics/internal/repository"
	"github.com/mikhailpashkov/metrics/internal/service"
	"github.com/mikhailpashkov/metrics/internal/utils"
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
	})

	logger.Debug("params read",
		"serverAddr", addr,
		"storeInterval", storeInterval,
		"fileStoragePath", fileStoragePath,
		"restore", restore,
	)

	// Dependencies ///////////////////
	logger.Debug("init dependencies")
	metricsRepository := repository.NewMetricsMemoryRepository()
	backupRepository := repository.NewFileBackupRepository(fileStoragePath)
	metricsService := service.NewMetricsService(metricsRepository, backupRepository)

	// Backup /////////////////////////
	if restore {
		logger.Debug("restore from backup")
		err := metricsService.Restore(context.Background())
		if err != nil {
			panic(err)
		}
	}
	logger.Debug("setup runtime backup")
	err := metricsService.SetupBackup(context.Background(), storeInterval)
	if err != nil {
		panic(err)
	}

	// Server /////////////////////////
	logger.Debug("setup server")

	r := chi.NewRouter()

	// наверняка, хорошая идея - использовать github.com/go-chi/chi/v5/middleware,
	// но в учебных целях используем самодельные
	r.Use(middleware.WithLogging(logger))
	r.Use(middleware.WithGZIPSupport(logger))

	r.Handle("/", handler.NewMetricsRootHandler(logger, metricsService))

	r.Handle("/value", handler.NewGetMetricsHandler(logger, metricsService))
	r.Handle("/value/{type}/{name}", handler.NewGetMetricsPathParamsHandler(logger, metricsService))

	r.Handle("/update", handler.NewUpdateMetricsHandler(logger, metricsService))
	r.Handle("/update/{type}/{name}/{value}", handler.NewUpdateMetricsPathParamsHandler(logger, metricsService))

	err = http.ListenAndServe(addr, r)
	if err != nil {
		logger.Error("failed to start server", "error", err)
		panic(err)
	}
}

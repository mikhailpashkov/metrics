package main

import (
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
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	logger.Info("SERVER")

	var addr string

	utils.GetParams([]utils.Param{
		&utils.StringParam{
			EnvName:       "ADDRESS",
			FlagName:      "a",
			FlagUsage:     "Server address",
			Default:       "localhost:8080",
			ValueConsumer: func(v string) { addr = v },
		},
	})

	logger.Debug("params read", "serverAddr", addr)

	metricsRepository := repository.NewMetricsMemoryRepository()
	metricsService := service.NewMetricsService(metricsRepository)

	handlers := []handler.MHandler{
		middleware.WithLogging(handler.NewGetMetricsHandler(logger, metricsService)),
		middleware.WithLogging(handler.NewGetListMetricsHandler(logger, metricsService)),
		middleware.WithLogging(handler.NewUpdateMetricsHandler(logger, metricsService)),
	}

	r := chi.NewRouter()
	for _, h := range handlers {
		r.Handle(h.GetUrlPattern(), h)
	}

	err := http.ListenAndServe(addr, r)
	if err != nil {
		logger.Error("failed to start server", "error", err)
		panic(err)
	}
}

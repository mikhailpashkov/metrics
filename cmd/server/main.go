package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikhailpashkov/metrics/internal/handler"
	"github.com/mikhailpashkov/metrics/internal/handler/middleware"
	"github.com/mikhailpashkov/metrics/internal/repository"
	"github.com/mikhailpashkov/metrics/internal/service"
	"github.com/mikhailpashkov/metrics/internal/utils"
	"go.uber.org/zap"
)

func main() {
	addr := utils.GetStringParam("ADDRESS", "a", "HTTP server address", "localhost:8080")
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Debug("Server addr:", zap.String("server addr", *addr))

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

	err = http.ListenAndServe(*addr, r)
	if err != nil {
		panic(err)
	}
}

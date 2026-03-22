package main

import (
	"net/http"

	"github.com/mikhailpashkov/metrics/internal/handler"
	"github.com/mikhailpashkov/metrics/internal/repository"
	"github.com/mikhailpashkov/metrics/internal/service"
)

const Addr = ":8080"

func main() {
	metricsRepository := repository.NewMetricsMemoryStorage()
	metricsService := service.NewMetrics(metricsRepository)
	metricsHandler := handler.NewMetrics(metricsService)

	mux := http.NewServeMux()
	mux.Handle(metricsHandler.GetUrlPattern(), metricsHandler)

	err := http.ListenAndServe(Addr, mux)
	if err != nil {
		panic(err)
	}
}

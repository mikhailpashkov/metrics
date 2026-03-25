package main

import (
	"net/http"

	"github.com/mikhailpashkov/metrics/internal/handler"
	"github.com/mikhailpashkov/metrics/internal/repository"
	"github.com/mikhailpashkov/metrics/internal/service"
)

const Addr = ":8080"

func main() {
	metricsRepository := repository.NewMetricsMemoryRepository()
	metricsService := service.NewMetricsService(metricsRepository)

	handlers := []handler.MHandler{
		handler.NewGetMetricsHandler(metricsService),
		handler.NewUpdateMetricsHandler(metricsService),
	}

	mux := http.NewServeMux()
	for _, h := range handlers {
		mux.Handle(h.GetUrlPattern(), h)
	}

	err := http.ListenAndServe(Addr, mux)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikhailpashkov/metrics/internal/handler"
	"github.com/mikhailpashkov/metrics/internal/repository"
	"github.com/mikhailpashkov/metrics/internal/service"
)

const Addr = ":8080"

func main() {
	fmt.Println("SERVER")
	metricsRepository := repository.NewMetricsMemoryRepository()
	metricsService := service.NewMetricsService(metricsRepository)

	handlers := []handler.MHandler{
		handler.NewGetMetricsHandler(metricsService),
		handler.NewGetListMetricsHandler(metricsService),
		handler.NewUpdateMetricsHandler(metricsService),
	}

	r := chi.NewRouter()
	for _, h := range handlers {
		r.Handle(h.GetUrlPattern(), h)
	}

	err := http.ListenAndServe(Addr, r)
	if err != nil {
		panic(err)
	}
}

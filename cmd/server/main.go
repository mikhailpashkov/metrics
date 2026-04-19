package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikhailpashkov/metrics/internal/handler"
	"github.com/mikhailpashkov/metrics/internal/repository"
	"github.com/mikhailpashkov/metrics/internal/service"
	"github.com/mikhailpashkov/metrics/internal/utils"
)

func main() {
	addr := utils.GetStringParam("ADDRESS", "a", "HTTP server address", "localhost:8080")
	flag.Parse()

	fmt.Println("SERVER", *addr)

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

	err := http.ListenAndServe(*addr, r)
	if err != nil {
		panic(err)
	}
}

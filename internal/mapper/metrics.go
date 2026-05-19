package mapper

import (
	"time"

	"github.com/mikhailpashkov/metrics/internal/dto"
	models "github.com/mikhailpashkov/metrics/internal/model"
)

func MetricsToGetMetricsResponse(metrics *models.Metrics) *dto.GetMetricsResponse {
	return &dto.GetMetricsResponse{
		ID:    metrics.Name,
		Type:  metrics.Type,
		Value: metrics.Value,
		Delta: metrics.Delta,
	}
}

func MetricsToUpdateMetricsRequest(metrics *models.Metrics) *dto.UpdateMetricsRequest {
	return &dto.UpdateMetricsRequest{
		ID:    metrics.Name,
		Type:  metrics.Type,
		Value: metrics.Value,
		Delta: metrics.Delta,
	}
}

func MetricsFromUpdateMetricsRequest(request dto.UpdateMetricsRequest) *models.Metrics {
	return &models.Metrics{
		ID:    -1,
		Type:  request.Type,
		Name:  request.ID,
		Delta: request.Delta,
		Value: request.Value,
		TS:    time.Now().UnixMilli(),
	}
}

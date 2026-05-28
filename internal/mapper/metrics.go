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

func MetricsFromMetricsDto(dto dto.MetricsDto) *models.Metrics {
	return &models.Metrics{
		ID:    models.MetricsNewID,
		Type:  dto.Type,
		Name:  dto.ID,
		Delta: dto.Delta,
		Value: dto.Value,
		TS:    time.Now().UnixMilli(),
	}
}

func MetricsToMetricsDto(m *models.Metrics) dto.MetricsDto {
	return dto.MetricsDto{
		Type:  m.Type,
		ID:    m.Name,
		Delta: m.Delta,
		Value: m.Value,
	}
}

func MetricsFromUpdateMetricsRequest(request dto.UpdateMetricsRequest) *models.Metrics {
	return MetricsFromMetricsDto(dto.MetricsDto(request))
}

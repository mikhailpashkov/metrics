package mapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mikhailpashkov/metrics/internal/dto"
	models "github.com/mikhailpashkov/metrics/internal/model"
)

func MetricsToGetMetricsResponse(metrics *models.Metrics) (*dto.GetMetricsResponse, error) {
	var value json.Number

	switch metrics.Type {
	case models.Gauge:
		formatFloat := strconv.FormatFloat(*metrics.Value, 'f', -1, 64)
		value = json.Number(formatFloat)
	case models.Counter:
		formatInt := strconv.FormatInt(*metrics.Delta, 10)
		value = json.Number(formatInt)
	default:
		return nil, fmt.Errorf("failed to map metrics. invalid metric type: %s", metrics.Type)
	}

	return &dto.GetMetricsResponse{
		ID:    metrics.Name,
		Type:  metrics.Type,
		Value: value,
	}, nil

}

func MetricsFromUpdateMetricsRequest(request dto.UpdateMetricsRequest) (*models.Metrics, error) {
	metrics := &models.Metrics{
		ID:    -1,
		Type:  request.Type,
		Name:  request.ID,
		Delta: nil,
		Value: nil,
		TS:    time.Now().UnixMilli(),
	}

	switch request.Type {
	case models.Counter:
		parseInt, err := request.Value.Int64()
		if err != nil {
			detailsErr := fmt.Errorf(
				"failed to map metrics. value conversion error. type: %s; id: %s; value: %s; err: %s",
				request.Type, request.ID, request.Value, err,
			)
			return nil, errors.Join(detailsErr, err)
		}
		metrics.Delta = &parseInt
	case models.Gauge:
		parseFloat, err := request.Value.Float64()
		if err != nil {
			detailsErr := fmt.Errorf(
				"failed to map metrics. value conversion error. type: %s; id: %s; value: %s; err: %s",
				request.Type, request.ID, request.Value, err,
			)
			return nil, errors.Join(detailsErr, err)
		}
		metrics.Value = &parseFloat

	default:
		return nil, fmt.Errorf("failed to map metrics. invalid metric type: %s", metrics.Type)
	}

	return metrics, nil
}

package models

func IsValidMetrics(metrics *Metrics) bool {
	switch metrics.Type {
	case Counter:
		return metrics.Delta != nil
	case Gauge:
		return metrics.Value != nil
	default:
		return false
	}
}

package service


import (
	"time"
)

type MetricsService struct {
	metrics *MetricsResponse
}

func newMetricsApi() *MetricsService {
	return &MetricsService{
		metrics: &MetricsResponse{
				TotalRequests: 0,
				AverageLatency: 0,
				LastRequestTime: time.Now(),
			},
		}
}

func (s *MetricsService) GetMetrics() (*MetricsResponse, error) {
 	return s.metrics, nil
}

func (s *MetricsService) UpdateMetrics(latency time.Duration) error {
	s.metrics.TotalRequests++
	s.metrics.AverageLatency = ((s.metrics.AverageLatency * time.Duration(s.metrics.TotalRequests-1)) + latency) / time.Duration(s.metrics.TotalRequests)
	s.metrics.LastRequestTime = time.Now()
	return nil
}


package service

import (
	"time"
)

type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type MetricsResponse struct {
	TotalRequests int       `json:"total_requests"`	
	AverageLatency time.Duration `json:"average_latency"`
	LastRequestTime time.Time `json:"last_request_time"`
}
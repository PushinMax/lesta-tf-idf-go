package handler

import (

	"github.com/PushinMax/lesta-tf-idf-go/internal/schema"
)
func ceilDiv(a, b int) int {
	if b == 0 {
		return 0
	}
	return (a + b - 1) / b
}


type LoginRequest struct {
	Login string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	Token string `json:"token" binding:"required"`	
}


type ErrorResponse struct {
    Message string `json:"error"`
}

// StatusResponse represents a generic status message response
type StatusResponse struct {
    Status string `json:"status"`
}

type VersionResponse struct {
    Version string `json:"version"`
}

// UploadResponse represents the response for file upload endpoints
type FileResponse struct {
    File string `json:"file"`
}

// TokenPairResponse represents the response for login and refresh endpoints
type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MetricsResponse struct {
	TotalRequests int       `json:"total_requests"`	
	AverageLatency string `json:"average_latency" example:"123.45ms"`
	LastRequestTime string `json:"last_request_time" example:"2025-06-16T10:00:00Z" format:"date-time"`
}


type CollectionListResponse struct {
    List []string `json:"list"`
}

// CollectionDocumentsResponse represents the response for collection documents endpoint
type CollectionDocumentsResponse struct {
    Documents []string `json:"documents"`
}

// CreateCollectionRequest represents request body for collection creation
type CreateCollectionRequest struct {
    Name string `json:"name" binding:"required"`
}

// DocumentListResponse represents the response for document list endpoint
type DocumentListResponse struct {
    List []string `json:"list"`
}

// DocumentResponse represents single document response
type DocumentResponse struct {
    File string `json:"file"`
}

// HuffmanResponse represents Huffman encoding response
type HuffmanResponse struct {
    Encoding    map[string]string `json:"encoding"`
    Compressed  string            `json:"compressed"`
}

type StatsResponse struct {
	Stat []schema.WordStat `json:"stat"`
}

type PageDataResponse struct {
	Words []schema.WordStat  `json:"words"`
	Total int  `json:"total"`
}

type UploadResponse struct {
	SessionID string            `json:"session_id"`	
	Page      int               `json:"page"`
	Words     []schema.WordStat `json:"words"`
	Total     int               `json:"total"`
}
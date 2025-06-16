package handler

import (
	"github.com/gin-gonic/gin"
 	"net/http"
 	"fmt"
	"time"
)

// @Summary Get API metrics
// @Description Get API usage metrics including total requests and latency
// @Tags metrics
// @Produce json
// @Success 200 {object} MetricsResponse
// @Failure 500 {object} ErrorResponse
// @Router /metrics [get]
func (h *Handler) getMetrics(c *gin.Context) {
	metrics, err := h.services.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
		"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	response := MetricsResponse{
        TotalRequests:   metrics.TotalRequests,
        AverageLatency:  metrics.AverageLatency.String(), 
        LastRequestTime: metrics.LastRequestTime.Format(time.RFC3339), 
    }
	c.JSON(http.StatusOK, response)
}



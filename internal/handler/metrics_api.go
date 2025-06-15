package handler

import (
	"github.com/gin-gonic/gin"
 	"net/http"
 	"fmt"
)

func (h *Handler) getMetrics(c *gin.Context) {
	metrics, err := h.services.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
		"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, metrics)
}



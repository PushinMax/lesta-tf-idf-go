package handler

import (
	"github.com/gin-gonic/gin"
)

// @Summary Get API version
// @Description Get current version of the API
// @Tags status
// @Produce json
// @Success 200 {object} VersionResponse
// @Failure 500 {object} ErrorResponse
// @Router /version [get]
func (h *Handler) version(c *gin.Context) {
	version, err := h.services.Version()
	if err != nil {
		c.JSON(500, gin.H{
		"error": err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"version": version,
	})
}

// @Summary Get API status
// @Description Check if API is running properly
// @Tags status
// @Produce json
// @Success 200 {object} StatusResponse
// @Failure 500 {object} ErrorResponse
// @Router /status [get]
func (h *Handler) status(c *gin.Context) {
	err := h.services.Status()
	if err != nil {
		c.JSON(500, gin.H{
		"error": err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"status": "OK",
	})
}
package handler

import (
	"github.com/gin-gonic/gin"
)

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
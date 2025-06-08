package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

func (h *Handler) login(c *gin.Context) {
	var request LoginRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %s", err.Error())})
		return
	}

	clientIP := c.GetHeader("X-Real-IP")
	if clientIP == "" {
		clientIP = c.ClientIP()
	}

	tokens, err := h.services.Login(request.Login, request.Password, clientIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) register(c *gin.Context) {
	var request LoginRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %s", err.Error())})
		return
	}

	err := h.services.Register(request.Login, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "You have successfully registered.",
	})
}
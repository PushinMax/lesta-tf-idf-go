package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)


// @Summary User Login
// @Tags auth
// @Description Authenticate user and get JWT tokens
// @Accept  json
// @Produce  json
// @Param input body LoginRequest true "Login credentials"
// @Success 200 {object} TokenPairResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
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

	c.JSON(http.StatusOK, TokenPairResponse{tokens.AccessToken, tokens.RefreshToken})
}


// @Summary User Registration
// @Tags auth
// @Description Register new user
// @Accept  json
// @Produce  json
// @Param input body LoginRequest true "Registration Data"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
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

// @Summary User Logout
// @Tags auth
// @Description Logout user and invalidate tokens
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} StatusResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /logout [get]
func (h *Handler) logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Invalid request body",
		})
		return
	}

	err := h.services.Logout(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, 
			ErrorResponse{
				Message: fmt.Sprintf("Failed: %s", err.Error()),
			},
		)
		return
	}

	c.JSON(http.StatusOK, StatusResponse{})
}


// @Summary User Refresh Token
// @Tags auth
// @Description Logout user and invalidate tokens
// @Accept  json
// @Produce  json
// @Param input body RefreshRequest true "Login credentials"
// @Success 200 {object} TokenPairResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login/refresh [POST]
func (h *Handler) refreshToken(c *gin.Context) {
	var request RefreshRequest
	if err := c.BindJSON(&request); err != nil {	
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %s", err.Error())})
	 	return
	}

	clientIP := c.GetHeader("X-Real-IP")
	if clientIP == "" {
		clientIP = c.ClientIP()
	}

	tokens, err := h.services.RefreshToken(request.Token, clientIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, tokens)
}
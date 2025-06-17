package handler

import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
)

// @Summary Change password
// @Description Change user's password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "New password"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user [patch]
func (h *Handler) changePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}
	var request ChangePasswordRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %s", err.Error())})
		return
	}
	err := h.services.ChangePassword(userID.(string), request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "You have successfully changed password.",
	})
}

// @Summary Delete account
// @Description Delete user's account
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StatusResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user [delete]
func (h *Handler) deleteAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.services.DeleteUser(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "You have successfully deleted your account.",
	})
}
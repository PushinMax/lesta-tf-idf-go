package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"net/http"
)


// @Summary Upload document
// @Description Upload a new document for processing
// @Tags documents
// @Accept multipart/form-data
// @Produce html
// @Security BearerAuth
// @Param Authorization header string true "Bearer {token}"
// @Param file formData file true "Document file"
// @Success 200 {object} UploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /upload [post]
func (h *Handler) handleUpload(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "File required"})
		return
	}

	response, err := h.services.UploadDocument(file, userID.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.HTML(200, "results.html", gin.H{
		"SessionID": response.SessionID,
		"Page":      response.Page,
		"Words":     response.Words,
		"Total":     response.Total,
	})

}

// @Summary Get page data
// @Description Get processed data for a specific page and session
// @Tags documents
// @Produce json
// @Param session path string true "Session ID"
// @Param page path int true "Page number"
// @Success 200 {object} PageDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /data/{session}/{page} [get]
func (h *Handler) getPageData(c * gin.Context) {
	sessionID := c.Param("session")
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.JSON(404, gin.H{"error": "page should be a number"})
		return
	}

	stat, err := h.services.GetPageData(sessionID, page)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"words": stat.Words,
		"total": stat.Total,
	})
}


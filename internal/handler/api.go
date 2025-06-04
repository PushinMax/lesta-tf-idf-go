package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
)
func (h *Handler) handleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "File required"})
		return
	}

	response, err := h.services.UploadFile(file)
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

func (h *Handler) getPageData(c * gin.Context) {
	sessionID := c.Param("session")
	page, _ := strconv.Atoi(c.Param("page"))

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
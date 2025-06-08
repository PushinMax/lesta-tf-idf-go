package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"fmt"
)

func (h *Handler) getListDocuments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}
	
	list, err := h.services.GetListFiles(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) getDocument(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}
	documentID := c.Param("document_id")
	content, err :=  h.services.GetFile(documentID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, struct{
		File string `json:"file"`
	}{
		File: content,
	})
}

func (h *Handler) getDocumentsStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}
	documentID := c.Param("document_id")
	list, err := h.services.GetFilesStats(documentID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, list)
}

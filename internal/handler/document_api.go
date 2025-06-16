package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"fmt"
)

// @Summary Get list of documents
// @Description Get all documents for authenticated user
// @Tags documents
// @Produce json
// @Security BearerAuth
// @Success 200 {object} DocumentListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /documents [get]
func (h *Handler) getListDocuments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}
	
	list, err := h.services.GetListDocuments(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, struct{
		List []string `json:"list"`
	}{
		List: list,
	})
}

// @Summary Get document content
// @Description Get content of specific document
// @Tags documents
// @Produce json
// @Security BearerAuth
// @Param document_id path string true "Document ID"
// @Success 200 {object} DocumentResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /documents/{document_id} [get]
func (h *Handler) getDocument(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}
	documentID := c.Param("document_id")
	content, err := h.services.GetDocument(documentID, userID.(string))
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

// @Summary Get document statistics
// @Description Get statistics for specific document
// @Tags documents
// @Produce json
// @Security BearerAuth
// @Param document_id path string true "Document ID"
// @Success 200 {object} StatsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /documents/{document_id}/statistics [get]
func (h *Handler) getDocumentsStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}
	documentID := c.Param("document_id")
	list, err := h.services.GetDocumentStats(documentID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, list)
}

// @Summary Delete document
// @Description Delete specific document
// @Tags documents
// @Security BearerAuth
// @Param document_id path string true "Document ID"
// @Success 200 {object} StatusResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /documents/{document_id} [delete]
func (h *Handler) deleteDocument(c * gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}
	documentID := c.Param("document_id")
	err := h.services.DeleteDocument(documentID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// @Summary Get Huffman encoding
// @Description Get Huffman encoding for specific document
// @Tags documents
// @Produce json
// @Security BearerAuth
// @Param document_id path string true "Document ID"
// @Success 200 {object} HuffmanResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /documents/{document_id}/haffman [get]
func (h *Handler) getHuffman(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}
	documentID := c.Param("document_id")
	huffman, err := h.services.GetHuffman(documentID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, huffman)
}


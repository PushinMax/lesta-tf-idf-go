package handler


import (
	"github.com/gin-gonic/gin"
 	"net/http"
 
 	"fmt"
)

func (h *Handler) getListCollections(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
	return
	}

	list, err := h.services.GetListCollections(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
		"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, struct {
		List []string `json:"list"`
		}{
		List: list,
	})
}
func (h *Handler) createCollection(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.services.CreateCollection(userID.(string), req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	c.Status(http.StatusOK)
}
func (h *Handler) getCollection(c *gin.Context) {
	 userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}

	collectionName := c.Param("collection_id")
		documents, err := h.services.GetDocumentsInCollection(userID.(string), collectionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed: %s", err.Error()),
	  })
		return
	}
	c.JSON(http.StatusOK, struct {
		Documents []string `json:"documents"`
		}{
		Documents: documents,
	})
}	
func (h *Handler) deleteCollection(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}

	collectionName := c.Param("collection_id")
	if err := h.services.DeleteCollection(userID.(string), collectionName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
		"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) getCollectionStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}

	collectionName := c.Param("collection_id")
	stats, err := h.services.GetCollectionStats(userID.(string), collectionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
		"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) addDocumentToCollection(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}

	collectionName := c.Param("collection_id")
	documentID := c.Param("document_id")

	if err := h.services.AddDocumentToCollection(userID.(string), collectionName, documentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
		"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) deleteDocumentFromCollection(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}

	collectionName := c.Param("collection_id")
	documentID := c.Param("document_id")

	if err := h.services.DeleteDocumentFromCollection(userID.(string), collectionName, documentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
		"error": fmt.Sprintf("Failed: %s", err.Error()),
		})
		return
	}
	c.Status(http.StatusOK)
}

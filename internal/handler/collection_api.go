package handler


import (
	"github.com/gin-gonic/gin"
 	"net/http"
 
 	"fmt"
)

// @Summary Get list of collections
// @Description Get all collections for authenticated user
// @Tags collections
// @Produce json
// @Security BearerAuth
// @Success 200 {object} CollectionListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /collections [get]
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
	c.JSON(http.StatusOK, CollectionListResponse{
		List: list,
	})
}

// @Summary Create new collection
// @Description Create a new collection for authenticated user
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCollectionRequest true "Collection creation request"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /collections/create [post]
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

// @Summary Get collection details
// @Description Get documents in specific collection
// @Tags collections
// @Produce json
// @Security BearerAuth
// @Param collection_name path string true "Collection ID"
// @Success 200 {object} CollectionDocumentsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /collections/{collection_name} [get]
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
	c.JSON(http.StatusOK, CollectionDocumentsResponse{
		Documents: documents,
	})
}	

// @Summary Delete collection
// @Description Delete specific collection
// @Tags collections
// @Security BearerAuth
// @Param collection_name path string true "Collection ID"
// @Success 200
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /collections/{collection_name} [delete]
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

// @Summary Get collection statistics
// @Description Get statistics for specific collection
// @Tags collections
// @Produce json
// @Security BearerAuth
// @Param collection_name path string true "Collection ID"
// @Success 200 {object} StatsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /collections/{collection_name}/statistics [get]
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

// @Summary Add document to collection
// @Description Add specific document to collection
// @Tags collections
// @Security BearerAuth
// @Param collection_id path string true "Collection ID"
// @Param document_name path string true "Document ID"
// @Success 200
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /collections/{collection_name}/{document_id} [post]
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

// @Summary Delete document from collection
// @Description Remove specific document from collection
// @Tags collections
// @Security BearerAuth
// @Param collection_id path string true "Collection ID"
// @Param document_id path string true "Document ID"
// @Success 200
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /collections/{collection_id}/{document_id} [delete]
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

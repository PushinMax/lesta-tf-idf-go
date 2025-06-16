package handler

import (
	"html/template"

	"github.com/PushinMax/lesta-tf-idf-go/internal/service"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func New(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	router.SetFuncMap(template.FuncMap{
		"ceilDiv": ceilDiv,
	})
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.Use(h.MetricsMiddleware())
	router.POST("/upload", h.JWTAuth(), h.handleUpload)
	router.GET("/data/:session/:page", h.getPageData)


	router.GET("/status", h.status)
	router.GET("/version", h.version)
	router.GET("/metrics", h.getMetrics)

	documents := router.Group("/documents")
	{
		documents.Use(h.JWTAuth())
		documents.GET("/", h.getListDocuments)
		documents.GET("/:document_id", h.getDocument)
		documents.GET("/:document_id/statistics", h.getDocumentsStats)
		documents.DELETE("/:document_id", h.deleteDocument)
		documents.GET("/:document_id/haffman", h.getHuffman)
	}
	collections := router.Group("/collections")
	{
		collections.Use(h.JWTAuth())
		collections.GET("/",  h.getListCollections)
		collections.POST("/create",  h.createCollection)
		collections.GET("/:collection_id",  h.getCollection)
		collections.GET("/:collection_id/statistics",  h.getCollectionStats)
		collections.DELETE("/:collection/:document_id",  h.deleteDocumentFromCollection)
		collections.POST("/:collection/:document_id",  h.addDocumentToCollection)
		collections.DELETE("/:collection/",  h.deleteCollection)
	}

	user := router.Group("/user")
	{
		user.Use(h.JWTAuth())
		user.PATCH("/:user_id", h.changePassword)
		user.DELETE("/:user_id", h.deleteAccount)
	}
	router.POST("/login", h.login)
	router.POST("/login/refresh", h.refreshToken)
	router.POST("/register", h.register)
	router.GET("/logout", h.JWTAuth(),  h.logout)

	return router
}

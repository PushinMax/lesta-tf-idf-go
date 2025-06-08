package handler

import (
	"html/template"

	"github.com/PushinMax/lesta-tf-idf-go/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func New(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	router.SetFuncMap(template.FuncMap{
		"ceilDiv": ceilDiv,
	})
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	
	router.POST("/upload", h.JWTAuth(), h.handleUpload)
	router.GET("/data/:session/:page", h.getPageData)


	router.GET("/status", h.status)
	router.GET("/version", h.version)
	router.GET("/metrics", nil)

	documents := router.Group("/documents")
	{
		documents.Use(h.JWTAuth())
		documents.GET("/", h.getListDocuments)
		documents.GET("/:document_id", h.getDocument)
		documents.GET("/:document_id/statistics", h.getDocumentsStats)
		documents.DELETE("/:document_id", nil)
	}
	collections := router.Group("/collections")
	{
		collections.Use(h.JWTAuth())
		collections.GET("/", nil)
		collections.GET("/:collection_id", nil)
		collections.GET("/:collection_id/statistics", nil)
		collections.DELETE("/:collection/:document_id", nil)
		collections.POST("/:collection/:document_id", nil)
	}

	user := router.Group("/user")
	{
		user.Use(h.JWTAuth())
		user.PATCH("/:user_id", h.changePassword)
		user.DELETE("/:user_id", nil)
	}
	router.POST("/login", h.login)
	router.POST("/register", h.register)
	router.GET("/logout", nil)

	return router
}

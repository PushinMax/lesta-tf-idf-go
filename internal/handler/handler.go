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

	router.POST("/upload", h.handleUpload)
	router.GET("/data/:session/:page", h.getPageData)


	router.GET("/status", h.status)
	router.GET("/version", h.version)


	return router
}

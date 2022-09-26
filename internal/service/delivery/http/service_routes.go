package http

import "github.com/gin-gonic/gin"

func NewServiceRoutes(group *gin.RouterGroup, h *serviceHandler) {
	group.GET("/:key", h.Get())
	group.POST("/", h.Set())
}

package sse

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	sseRouter := router.Group("/sse")

	sseRouter.GET("", sseHandler)
}

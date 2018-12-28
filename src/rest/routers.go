package rest

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {

	router.GET("/admin", func(context *gin.Context) {

	})

	router.POST("/api/generate", generate)
	router.NoRoute(dwzHandler)

}

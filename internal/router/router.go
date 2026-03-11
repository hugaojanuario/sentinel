package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/sentinel/internal/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/containers", controllers.ListContainers)

	router.POST("/containers/:id/restart", controllers.RestartContainer)

	router.GET("/containers/:id/logs", controllers.GetContainerLogs)

	router.GET("/containers/:id/stats", controllers.GetContainerStats)

	return router
}

package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/sentinel/internal/controllers"
)

func SetupRouter(auth controllers.AuthController) *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())

	router.POST("/login", auth.Login)

	router.GET("/containers", controllers.ListContainers)

	router.POST("/containers/:id/restart", controllers.RestartContainer)

	router.GET("/containers/:id/logs", controllers.GetContainerLogs)

	router.GET("/containers/:id/stats", controllers.GetContainerStats)

	router.GET("/containers/:id/logs/stream", controllers.StreamLogs)

	return router
}

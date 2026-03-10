package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/sentinel/internal/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	router.GET("/containers", controllers.ListContainers)

	router.POST("/containers/:id/restart", controllers.RestartContainer)

	router.GET("/containers/:id/logs", controllers.GetContainerLogs)

	return router
}

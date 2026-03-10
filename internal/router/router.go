package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/sentinel/internal/controllers"
	"github.com/hugaojanuario/sentinel/internal/docker"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	router.GET("/containers", controllers.ListContainers)

	router.POST("/containers/:id/restart", func(c *gin.Context) {
		id := c.Param("id")
		err := docker.RestartContainer(id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "restarted container"})
	})

	router.GET("/containers/:id/logs", func(c *gin.Context){
		id := c.Param("id")
		logs, err := docker.GetContainerLogs(id)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return
		}

		c.String(http.StatusOK, logs)

	})

	return router
}

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/sentinel/internal/services"
)

func ListContainers(c *gin.Context) {
	containers, err := services.ListContainers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, containers)

}

func RestartContainer(c *gin.Context) {

	id := c.Param("id")
	err := services.RestartContainer(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "container restarted"})
}

func GetContainerLogs(c *gin.Context) {
	id := c.Param("id")
	logs, err := services.GetContainerLogs(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, logs)

}

func GetContainerStats(c *gin.Context) {

	id := c.Param("id")

	stats, err := services.GetContainerStats(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

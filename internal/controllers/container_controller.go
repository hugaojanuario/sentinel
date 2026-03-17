package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/sentinel/internal/docker"
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

func StreamLogs(c *gin.Context) {
	id := c.Param("id")

	reader, err := docker.StreamContainerLogs(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer reader.Close()

	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	buf := make([]byte, 1024)

	for {
		n, err := reader.Read(buf)
		if err != nil {
			break
		}

		c.Writer.Write(buf[:n])
		c.Writer.Flush()
	}
}

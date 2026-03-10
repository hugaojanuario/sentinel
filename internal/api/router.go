package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/sentinel/internal/docker"
)

func NewRouter ()*gin.Engine{
	router := gin.Default()

	router.GET("/status", func (c *gin.Context){
		c.JSON(http.StatusOK, gin.H{"message":"OK"})
	})

	router.GET("/containers", func(c *gin.Context){
		containers, err := docker.ListContainers()

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return
		}

		c.JSON(http.StatusOK, containers)
	})

	return router
}
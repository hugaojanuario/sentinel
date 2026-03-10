package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter ()*gin.Engine{
	router := gin.Default()

	router.GET("/status", func (c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message":"OK",
		})
	})

	return router
}
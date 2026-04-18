package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/sentinel/internal/controllers"
	"github.com/hugaojanuario/sentinel/internal/middleware"
)

func SetupRouter(auth *controllers.AuthController) *gin.Engine {
	//PUBLIC
	public := gin.Default()
	public.Use(cors.Default())

	public.POST("/register", auth.Register)
	public.POST("/login", auth.Login)

	//PROTECTED
	protected := public.Group("/containers", middleware.AuthMiddleware())

	protected.GET("/", controllers.ListContainers)
	protected.POST("/:id/restart", controllers.RestartContainer)
	protected.GET("/:id/logs", controllers.GetContainerLogs)
	protected.GET("/:id/stats", controllers.GetContainerStats)
	protected.GET("/:id/logs/stream", controllers.StreamLogs)

	return public
}

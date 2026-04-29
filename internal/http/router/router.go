package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	handler2 "github.com/hugaojanuario/sentinel/internal/http/handler"
	"github.com/hugaojanuario/sentinel/internal/http/middleware"
)

func SetupRouter(auth *handler2.AuthController) *gin.Engine {
	//PUBLIC
	public := gin.Default()
	public.Use(cors.Default())

	public.POST("/register", auth.Register)
	public.POST("/login", auth.Login)

	//PROTECTED
	protected := public.Group("/containers", middleware.AuthMiddleware())

	protected.GET("", handler2.ListContainers)
	protected.POST("/:id/restart", handler2.RestartContainer)
	protected.GET("/:id/logs", handler2.GetContainerLogs)
	protected.GET("/:id/stats", handler2.GetContainerStats)
	protected.GET("/:id/logs/stream", handler2.StreamLogs)

	return public
}

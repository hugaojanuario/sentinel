package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/sentinel/internal/http/handler"
	"github.com/hugaojanuario/sentinel/internal/http/middleware"
)

func SetupRouter(auth *handler.AuthController, secret []byte) *gin.Engine {
	//PUBLIC
	public := gin.Default()
	public.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	public.POST("/register", auth.Register)
	public.POST("/login", auth.Login)

	//PROTECTED
	protected := public.Group("/containers", middleware.AuthMiddleware(secret))

	protected.GET("", handler.ListContainers)
	protected.POST("/:id/restart", handler.RestartContainer)
	protected.GET("/:id/logs", handler.GetContainerLogs)
	protected.GET("/:id/stats", handler.GetContainerStats)
	protected.GET("/:id/logs/stream", handler.StreamLogs)

	return public
}

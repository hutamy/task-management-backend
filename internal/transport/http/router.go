package http

import (
	"net/http"
	"task-management-backend/internal/transport/http/handlers"
	"task-management-backend/middleware"

	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	Auth      *handlers.AuthHandler
	Task      *handlers.TaskHandler
	JwtSecret string
}

func RegisterRoutes(g *gin.Engine, deps RouterDeps) {
	g.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Task Management API"})
	})

	api := g.Group("/api")
	{
		api.POST("/login", deps.Auth.Login)
	}

	protected := api.Group("/tasks")
	protected.Use(middleware.JWTMiddleware(deps.JwtSecret))
	{
		protected.GET("", deps.Task.GetTasks)
		protected.POST("", deps.Task.CreateTask)
		protected.PUT("/:id", deps.Task.UpdateTask)
		protected.DELETE("/:id", deps.Task.DeleteTask)
	}
}

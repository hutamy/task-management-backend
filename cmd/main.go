package main

import (
	"fmt"
	"log"
	"task-management-backend/config"
	"task-management-backend/internal/cache"
	"task-management-backend/internal/repository"
	ht "task-management-backend/internal/transport/http"
	"task-management-backend/internal/transport/http/handlers"
	"task-management-backend/internal/usecase/auth"
	"task-management-backend/internal/usecase/task"
	"task-management-backend/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	db, err := config.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	defer db.Close()

	taskCache := cache.NewTaskCache(time.Duration(cfg.CacheDuration) * time.Hour)

	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)

	authUC := auth.NewAuthUseCase(cfg.JwtSecret, userRepo)
	taskUC := task.NewTaskUseCase(taskRepo, taskCache)

	authHandler := handlers.NewAuthHandler(authUC)
	taskHandler := handlers.NewTaskHandler(taskUC)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	ht.RegisterRoutes(router, ht.RouterDeps{
		Auth:      authHandler,
		Task:      taskHandler,
		JwtSecret: cfg.JwtSecret,
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

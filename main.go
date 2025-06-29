package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task-api/database"
	"task-api/handlers"
	"task-api/middleware"
	"task-api/repositories"
	"task-api/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	userRepo := repositories.NewUserRepository(database.DB)
	taskRepo := repositories.NewTaskRepository(database.DB)

	authService := services.NewAuthService(userRepo)
	taskService := services.NewTaskService(taskRepo)

	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "connected",
			"service":  "task-api",
			"version":  "1.0.0",
		})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.GET("/profile", middleware.AuthRequired(), authHandler.GetProfile)
		}

		tasks := v1.Group("/tasks")
		tasks.Use(middleware.AuthRequired())
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.GetUserTasks)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
			tasks.POST("/:id/complete", taskHandler.CompleteTask)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.AuthRequired())
		{
			admin.GET("/tasks", taskHandler.GetAllTasks)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
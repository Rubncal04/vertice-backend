package main

import (
	"log"
	"os"

	"vertice-backend/config"
	"vertice-backend/internal/domain"
	"vertice-backend/internal/handler"
	middleware2 "vertice-backend/internal/middleware"
	"vertice-backend/internal/repository"
	"vertice-backend/internal/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = godotenv.Load()

	config.InitDB()
	if err := config.DB.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("Error in migration: %v", err)
	}

	repo := repository.NewUserGormRepository()
	service := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(service)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api/v1")

	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)

	userGroup := api.Group("/users")
	userGroup.Use(middleware2.JWTMiddleware())
	userGroup.GET("/profile", userHandler.Profile)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}

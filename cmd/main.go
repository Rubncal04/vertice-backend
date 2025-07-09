package main

import (
	"log"
	"os"

	"vertice-backend/config"
	"vertice-backend/internal/repository"
	"vertice-backend/internal/service"
	"vertice-backend/migrations"
	"vertice-backend/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = godotenv.Load()

	config.InitDB()
	if err := migrations.AutoMigrateAll(config.DB); err != nil {
		log.Fatalf("Error in migration: %v", err)
	}

	userRepo := repository.NewUserGormRepository()
	productRepo := repository.NewProductGormRepository()

	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.RegisterAllRoutes(e, userService, productService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}

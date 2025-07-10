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

	_ "vertice-backend/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Vertice Backend API
// @version 1.0
// @description API System for managing products and orders of Vertice
// @termsOfService http://swagger.io/terms/

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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

	orderRepo := repository.NewOrderGormRepository()
	orderService := service.NewOrderService(orderRepo, productRepo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	routesDependencies := routes.AppDependencies{
		UserService:    userService,
		ProductService: productService,
		OrderService:   orderService,
	}

	routes.RegisterAllRoutes(e, routesDependencies)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}

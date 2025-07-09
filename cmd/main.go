package main

import (
	"log"
	"os"

	"vertice-backend/config"
	"vertice-backend/internal/handler"
	middleware2 "vertice-backend/internal/middleware"
	"vertice-backend/internal/repository"
	"vertice-backend/internal/service"
	"vertice-backend/migrations"

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

	repo := repository.NewUserGormRepository()
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)

	productRepo := repository.NewProductGormRepository()
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api/v1")

	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)

	userGroup := api.Group("/users")
	userGroup.Use(middleware2.JWTMiddleware())
	userGroup.GET("/profile", userHandler.Profile)

	productGroup := api.Group("/products")
	productGroup.Use(middleware2.JWTMiddleware())
	productGroup.POST("", productHandler.CreateProduct)
	productGroup.GET("", productHandler.ListProducts)
	productGroup.GET("/:id", productHandler.GetProduct)
	productGroup.PUT("/:id", productHandler.UpdateProduct)
	productGroup.DELETE("/:id", productHandler.DeleteProduct)
	productGroup.PATCH("/:id/stock", productHandler.UpdateProductStock)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}

package routes

import (
	"vertice-backend/internal/handler"
	"vertice-backend/internal/middleware"
	"vertice-backend/internal/service"

	"github.com/labstack/echo/v4"
)

func RegisterProductRoutes(e *echo.Echo, productService *service.ProductService) {
	productHandler := handler.NewProductHandler(productService)

	api := e.Group("/api/v1")

	products := api.Group("/products", middleware.JWTMiddleware())

	products.POST("", productHandler.CreateProduct)
	products.GET("", productHandler.ListProducts)
	products.GET("/:id", productHandler.GetProduct)
	products.PATCH("/:id", productHandler.UpdateProduct)
	products.DELETE("/:id", productHandler.DeleteProduct)
	products.PATCH("/:id/stock", productHandler.UpdateProductStock)
}

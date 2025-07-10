package routes

import (
	"vertice-backend/internal/handler"
	"vertice-backend/internal/middleware"
	"vertice-backend/internal/service"

	"github.com/labstack/echo/v4"
)

func RegisterOrderRoutes(e *echo.Echo, orderService *service.OrderService) {
	orderHandler := handler.NewOrderHandler(orderService)

	api := e.Group("/api/v1")
	orders := api.Group("/orders", middleware.JWTMiddleware())

	orders.POST("", orderHandler.CreateOrder)
	orders.GET("", orderHandler.ListOrders)
	orders.GET("/:id", orderHandler.GetOrder)
	orders.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
	orders.POST("/:id/cancel", orderHandler.CancelOrder)
	orders.DELETE("/:id", orderHandler.DeleteOrder)
}

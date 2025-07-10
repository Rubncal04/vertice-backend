package routes

import (
	"vertice-backend/internal/service"

	"github.com/labstack/echo/v4"
)

type AppDependencies struct {
	UserService    *service.UserService
	ProductService *service.ProductService
	OrderService   *service.OrderService
}

func RegisterAllRoutes(e *echo.Echo, deps AppDependencies) {
	RegisterUserRoutes(e, deps.UserService)
	RegisterProductRoutes(e, deps.ProductService)
	RegisterOrderRoutes(e, deps.OrderService)
}

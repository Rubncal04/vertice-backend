package routes

import (
	"vertice-backend/internal/service"

	"github.com/labstack/echo/v4"
)

func RegisterAllRoutes(e *echo.Echo, userService *service.UserService, productService *service.ProductService) {
	RegisterUserRoutes(e, userService)
	RegisterProductRoutes(e, productService)
}

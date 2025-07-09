package routes

import (
	"vertice-backend/internal/handler"
	"vertice-backend/internal/middleware"
	"vertice-backend/internal/service"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(e *echo.Echo, userService *service.UserService) {
	userHandler := handler.NewUserHandler(userService)

	api := e.Group("/api/v1")

	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)

	users := api.Group("/users", middleware.JWTMiddleware())

	users.GET("/profile", userHandler.Profile)
}

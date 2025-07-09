package handler

import (
	"context"
	"net/http"
	"vertice-backend/pkg"

	"vertice-backend/internal/domain"

	"github.com/labstack/echo/v4"
)

type UserServiceInterface interface {
	Register(ctx context.Context, name, email, password string) (*domain.User, error)
	Authenticate(ctx context.Context, email, password string) (*domain.User, error)
	GetProfile(ctx context.Context, id uint) (*domain.User, error)
}

type UserHandler struct {
	service UserServiceInterface
}

func NewUserHandler(service UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	user, err := h.service.Register(c.Request().Context(), req.Name, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"id": user.ID, "name": user.Name, "email": user.Email})
}

func (h *UserHandler) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	user, err := h.service.Authenticate(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}
	token, err := pkg.GenerateJWT(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not generate token"})
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (h *UserHandler) Profile(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
	}
	user, err := h.service.GetProfile(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
	}
	return c.JSON(http.StatusOK, echo.Map{"id": user.ID, "name": user.Name, "email": user.Email})
}

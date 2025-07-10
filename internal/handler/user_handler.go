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
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"password123"`
}

type loginRequest struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"password123"`
}

type userResponse struct {
	ID    uint   `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

type loginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body registerRequest true "User data"
// @Success 201 {object} userResponse
// @Failure 400 {object} ErrorResponse
// @Router /users/register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	user, err := h.service.Register(c.Request().Context(), req.Name, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, userResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

// Login godoc
// @Summary Login
// @Description Authenticate a user and return a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body loginRequest true "Login credentials"
// @Success 200 {object} loginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/login [post]
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
	return c.JSON(http.StatusOK, loginResponse{Token: token})
}

// Profile godoc
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/profile [get]
func (h *UserHandler) Profile(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
	}
	user, err := h.service.GetProfile(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
	}
	return c.JSON(http.StatusOK, userResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

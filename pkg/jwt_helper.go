package pkg

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func GetUserIDFromJWTContext(c echo.Context) (uint, error) {
	header := c.Request().Header.Get("Authorization")
	if header == "" {
		return 0, echo.ErrUnauthorized
	}
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, echo.ErrUnauthorized
	}
	claims, err := ParseJWT(parts[1])
	if err != nil {
		return 0, echo.ErrUnauthorized
	}
	return claims.UserID, nil
}

package handler

import (
	"net/http"
	"strconv"

	"vertice-backend/internal/service"
	"vertice-backend/pkg"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	type reqBody struct {
		Code        string  `json:"code"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}
	var body reqBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	product, err := h.service.CreateProduct(c.Request().Context(), userID, body.Code, body.Name, body.Description, body.Price, body.Stock)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) ListProducts(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	products, err := h.service.GetProductsByUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid product id")
	}
	product, err := h.service.GetProduct(c.Request().Context(), uint(id), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid product id")
	}
	type reqBody struct {
		Code        *string  `json:"code"`
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Price       *float64 `json:"price"`
		Stock       *int     `json:"stock"`
	}
	var body reqBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	product, err := h.service.UpdateProduct(c.Request().Context(), uint(id), userID, body.Code, body.Name, body.Description, body.Price, body.Stock)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProductStock(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid product id")
	}
	type reqBody struct {
		StockDelta int `json:"stockDelta"`
	}
	var body reqBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	product, err := h.service.UpdateProductStock(c.Request().Context(), uint(id), userID, body.StockDelta)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid product id")
	}
	if err := h.service.DeleteProduct(c.Request().Context(), uint(id), userID); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

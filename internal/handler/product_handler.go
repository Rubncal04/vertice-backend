package handler

import (
	"net/http"
	"strconv"

	"vertice-backend/internal/domain"
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

type createProductRequest struct {
	Code        string  `json:"code" example:"PROD001"`
	Name        string  `json:"name" example:"Laptop Gaming"`
	Description string  `json:"description" example:"Laptop para gaming de alta performance"`
	Price       float64 `json:"price" example:"1299.99"`
	Stock       int     `json:"stock" example:"10"`
}

type updateProductRequest struct {
	Code        *string  `json:"code,omitempty" example:"PROD001"`
	Name        *string  `json:"name,omitempty" example:"Laptop Gaming Pro"`
	Description *string  `json:"description,omitempty" example:"Laptop para gaming de alta performance actualizada"`
	Price       *float64 `json:"price,omitempty" example:"1399.99"`
	Stock       *int     `json:"stock,omitempty" example:"15"`
}

type updateStockRequest struct {
	StockDelta int `json:"stockDelta" example:"5"`
}

type ProductResponse struct {
	ID          uint    `json:"id" example:"1"`
	Code        string  `json:"code" example:"PROD001"`
	Name        string  `json:"name" example:"Laptop"`
	Description string  `json:"description" example:"Laptop para gaming"`
	Price       float64 `json:"price" example:"1299.99"`
	Stock       int     `json:"stock" example:"10"`
}

func toProductResponse(p *domain.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product for the authenticated user
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body createProductRequest true "Product data"
// @Success 201 {object} ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	var body createProductRequest
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	product, err := h.service.CreateProduct(c.Request().Context(), userID, body.Code, body.Name, body.Description, body.Price, body.Stock)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, toProductResponse(product))
}

// ListProducts godoc
// @Summary List products of the authenticated user
// @Description Get all products of the authenticated user
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ProductResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [get]
func (h *ProductHandler) ListProducts(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	products, err := h.service.GetProductsByUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	responses := make([]ProductResponse, len(products))
	for i, p := range products {
		responses[i] = toProductResponse(p)
	}
	return c.JSON(http.StatusOK, responses)
}

// GetProduct godoc
// @Summary Get a specific product
// @Description Get a specific product by the authenticated user's ID
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [get]
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
	return c.JSON(http.StatusOK, toProductResponse(product))
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product of the authenticated user
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param product body updateProductRequest true "Data to update"
// @Success 200 {object} ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [put]
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
	var body updateProductRequest
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	product, err := h.service.UpdateProduct(c.Request().Context(), uint(id), userID, body.Code, body.Name, body.Description, body.Price, body.Stock)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, toProductResponse(product))
}

// UpdateProductStock godoc
// @Summary Update product stock
// @Description Update the stock of a product (increment or decrement)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param stock body updateStockRequest true "Stock change data"
// @Success 200 {object} ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id}/stock [patch]
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
	var body updateStockRequest
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	product, err := h.service.UpdateProductStock(c.Request().Context(), uint(id), userID, body.StockDelta)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, toProductResponse(product))
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product of the authenticated user
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [delete]
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

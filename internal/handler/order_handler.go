package handler

import (
	"net/http"
	"strconv"
	"time"

	"vertice-backend/internal/domain"
	"vertice-backend/internal/service"
	"vertice-backend/pkg"

	"github.com/labstack/echo/v4"
)

type ProductSummary struct {
	ID    uint    `json:"id" example:"1"`
	Code  string  `json:"code" example:"PROD001"`
	Name  string  `json:"name" example:"Laptop Gaming"`
	Price float64 `json:"price" example:"1299.99"`
}

type OrderItemResponse struct {
	ID        uint           `json:"id" example:"1"`
	ProductID uint           `json:"product_id" example:"1"`
	Product   ProductSummary `json:"product"`
	Quantity  int            `json:"quantity" example:"2"`
	UnitPrice float64        `json:"unit_price" example:"1299.99"`
	Subtotal  float64        `json:"subtotal" example:"2599.98"`
}

type OrderResponse struct {
	ID          uint                `json:"id" example:"1"`
	Status      string              `json:"status" example:"pending"`
	TotalAmount float64             `json:"total_amount" example:"2599.98"`
	Items       []OrderItemResponse `json:"items"`
	CreatedAt   time.Time           `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt   time.Time           `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

type updateOrderStatusRequest struct {
	Status string `json:"status" example:"processing"`
}

func toProductSummary(p domain.Product) ProductSummary {
	return ProductSummary{
		ID:    p.ID,
		Code:  p.Code,
		Name:  p.Name,
		Price: p.Price,
	}
}

func toOrderItemResponse(item domain.OrderItem) OrderItemResponse {
	var prodSummary ProductSummary
	if item.Product.ID != 0 {
		prodSummary = toProductSummary(item.Product)
	}
	return OrderItemResponse{
		ID:        item.ID,
		ProductID: item.ProductID,
		Product:   prodSummary,
		Quantity:  item.Quantity,
		UnitPrice: item.UnitPrice,
		Subtotal:  item.Subtotal,
	}
}

func toOrderResponse(order *domain.Order) OrderResponse {
	items := make([]OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = toOrderItemResponse(item)
	}
	return OrderResponse{
		ID:          order.ID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		Items:       items,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body service.CreateOrderRequest true "Order data"
// @Success 201 {object} OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	var req service.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	order, err := h.service.CreateOrder(c.Request().Context(), userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, toOrderResponse(order))
}

// ListOrders godoc
// @Summary List orders of the authenticated user
// @Description Get all orders of the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} OrderResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [get]
func (h *OrderHandler) ListOrders(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	orders, err := h.service.GetOrdersByUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	resp := make([]OrderResponse, len(orders))
	for i, order := range orders {
		resp[i] = toOrderResponse(order)
	}
	return c.JSON(http.StatusOK, resp)
}

// GetOrder godoc
// @Summary Get a specific order
// @Description Get a specific order by the authenticated user's ID
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la orden"
// @Success 200 {object} OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order id")
	}
	order, err := h.service.GetOrder(c.Request().Context(), uint(id), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, toOrderResponse(order))
}

// UpdateOrderStatus godoc
// @Summary Update the status of an order
// @Description Update the status of an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la orden"
// @Param status body updateOrderStatusRequest true "Nuevo estado"
// @Success 200 {object} OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /orders/{id}/status [patch]
func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order id")
	}
	var body updateOrderStatusRequest
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	order, err := h.service.UpdateOrderStatus(c.Request().Context(), uint(id), userID, domain.OrderStatus(body.Status))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, toOrderResponse(order))
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel an existing order of the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la orden"
// @Success 200 {object} OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /orders/{id}/cancel [patch]
func (h *OrderHandler) CancelOrder(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order id")
	}
	order, err := h.service.CancelOrder(c.Request().Context(), uint(id), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, toOrderResponse(order))
}

// DeleteOrder godoc
// @Summary Delete an order
// @Description Delete an order of the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la orden"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(c echo.Context) error {
	userID, err := pkg.GetUserIDFromJWTContext(c)
	if err != nil {
		return err
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order id")
	}
	if err := h.service.DeleteOrder(c.Request().Context(), uint(id), userID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

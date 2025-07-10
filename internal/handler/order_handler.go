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
	ID    uint    `json:"id"`
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type OrderItemResponse struct {
	ID        uint           `json:"id"`
	ProductID uint           `json:"product_id"`
	Product   ProductSummary `json:"product"`
	Quantity  int            `json:"quantity"`
	UnitPrice float64        `json:"unit_price"`
	Subtotal  float64        `json:"subtotal"`
}

type OrderResponse struct {
	ID          uint                `json:"id"`
	Status      string              `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	Items       []OrderItemResponse `json:"items"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
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
	type reqBody struct {
		Status string `json:"status"`
	}
	var body reqBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	order, err := h.service.UpdateOrderStatus(c.Request().Context(), uint(id), userID, domain.OrderStatus(body.Status))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, toOrderResponse(order))
}

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

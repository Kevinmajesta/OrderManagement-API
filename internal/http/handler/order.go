package handler

import (
	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/internal/http/binder"
	"Kevinmajesta/OrderManagementAPI/internal/service"
	"Kevinmajesta/OrderManagementAPI/pkg/response"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var req binder.OrderCreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid request"))
	}

	var orderItems []entity.OrderItem
	for _, item := range req.Items {
		orderItems = append(orderItems, entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	order := &entity.Order{
		UserID:        req.UserID,
		PaymentMethod: req.PaymentMethod,
		PaidAmount:    req.PaidAmount,
		OrderItems:    orderItems,
	}

	if err := h.orderService.CreateOrder(order); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusCreated, response.SuccessResponse(http.StatusCreated, "Order created successfully", order))
}

func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
	orderIDParam := c.Param("order_id")
	orderID, err := uuid.Parse(orderIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid order_id"))
	}

	var req binder.OrderUpdateStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid request body"))
	}

	// isi manual order ID dari path param
	req.OrderID = orderID

	if err := h.orderService.UpdateOrderStatus(req.OrderID, req.Status); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Order status updated", nil))
}

func (h *OrderHandler) GetOrderHistory(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "user_id is required"))
	}

	orders, err := h.orderService.GetOrderHistory(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "order history fetched", orders))
}

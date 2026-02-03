package handler

import (
	"Kevinmajesta/OrderManagementAPI/internal/service"
	"Kevinmajesta/OrderManagementAPI/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MidtransHandler struct {
	orderService service.OrderService
}

func NewMidtransHandler(orderService service.OrderService) *MidtransHandler {
	return &MidtransHandler{orderService: orderService}
}

func (h *MidtransHandler) HandleNotification(c echo.Context) error {
	var notification map[string]interface{}

	if err := c.Bind(&notification); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid notification"))
	}

	orderID, ok := notification["order_id"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "order_id not found"))
	}

	transactionStatus, ok := notification["transaction_status"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "transaction_status not found"))
	}

	fraudStatus, _ := notification["fraud_status"].(string)

	// Determine order status based on Midtrans transaction status
	var orderStatus string
	switch transactionStatus {
	case "capture":
		if fraudStatus == "accept" {
			orderStatus = "paid"
		}
	case "settlement":
		orderStatus = "paid"
	case "pending":
		orderStatus = "pending"
	case "deny", "expire", "cancel":
		orderStatus = "cancelled"
	default:
		orderStatus = "pending"
	}

	// Update order status
	if err := h.orderService.UpdateOrderStatusByOrderID(orderID, orderStatus); err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

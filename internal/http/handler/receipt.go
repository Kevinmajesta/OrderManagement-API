package handler

import (
	"net/http"

	"Kevinmajesta/OrderManagementAPI/internal/http/binder"
	"Kevinmajesta/OrderManagementAPI/internal/service"
	"Kevinmajesta/OrderManagementAPI/pkg/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReceiptHandler struct {
	receiptService service.ReceiptService
}

func NewReceiptHandler(receiptService service.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{receiptService: receiptService}
}

func (h *ReceiptHandler) GenerateReceipt(c echo.Context) error {
	var req binder.ReceiptGenerateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid request"))
	}

	receipt, err := h.receiptService.GenerateReceipt(req.OrderID, req.UserID, req.CashierName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusCreated, response.SuccessResponse(http.StatusCreated, "receipt generated", receipt))
}

func (h *ReceiptHandler) GetReceiptByID(c echo.Context) error {
	receiptIDParam := c.Param("receipt_id")
	receiptID, err := uuid.Parse(receiptIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid receipt_id"))
	}

	receipt, err := h.receiptService.GetReceiptByID(receiptID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "receipt not found"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "receipt fetched", receipt))
}

func (h *ReceiptHandler) GetReceiptByOrderID(c echo.Context) error {
	orderIDParam := c.QueryParam("order_id")
	if orderIDParam == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "order_id is required"))
	}

	orderID, err := uuid.Parse(orderIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid order_id"))
	}

	receipt, err := h.receiptService.GetReceiptByOrderID(orderID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "receipt not found"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "receipt fetched", receipt))
}

func (h *ReceiptHandler) GetReceiptsByUserID(c echo.Context) error {
	userIDParam := c.QueryParam("user_id")
	if userIDParam == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "user_id is required"))
	}

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid user_id"))
	}

	receipts, err := h.receiptService.GetReceiptsByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "receipts fetched", receipts))
}

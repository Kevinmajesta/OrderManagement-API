package handler

import (
	"net/http"

	"Kevinmajesta/OrderManagementAPI/internal/http/binder"
	"Kevinmajesta/OrderManagementAPI/internal/service"
	"Kevinmajesta/OrderManagementAPI/pkg/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) AddItem(c echo.Context) error {
	var req binder.CartAddItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid request"))
	}

	cart, err := h.cartService.AddItem(req.UserID, req.ProductID, req.Quantity)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "item added", cart))
}

func (h *CartHandler) UpdateItem(c echo.Context) error {
	cartItemIDParam := c.Param("cart_item_id")
	cartItemID, err := uuid.Parse(cartItemIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid cart_item_id"))
	}

	var req binder.CartUpdateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid request"))
	}

	cart, err := h.cartService.UpdateItem(cartItemID, req.Quantity)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "item updated", cart))
}

func (h *CartHandler) RemoveItem(c echo.Context) error {
	cartItemIDParam := c.Param("cart_item_id")
	cartItemID, err := uuid.Parse(cartItemIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid cart_item_id"))
	}

	if err := h.cartService.RemoveItem(cartItemID); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "item removed", nil))
}

func (h *CartHandler) GetCart(c echo.Context) error {
	userIDParam := c.QueryParam("user_id")
	if userIDParam == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "user_id is required"))
	}

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid user_id"))
	}

	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "cart fetched", cart))
}

func (h *CartHandler) Checkout(c echo.Context) error {
	var req binder.CartCheckoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid request"))
	}

	order, err := h.cartService.Checkout(req.UserID, req.PaymentMethod, req.PaidAmount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusCreated, response.SuccessResponse(http.StatusCreated, "checkout success", order))
}

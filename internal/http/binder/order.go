package binder

import "github.com/google/uuid"

type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

type OrderCreateRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required"` // <— wajib ada ini
	PaymentMethod string    `json:"payment_method"`
	PaidAmount    float64   `json:"paid_amount"`
	Items         []struct {
		ProductID uuid.UUID `json:"product_id"`
		Quantity  int       `json:"quantity"`
	} `json:"items"`
}

type OrderUpdateStatusRequest struct {
	OrderID uuid.UUID `param:"order_id" json:"order_id" validate:"required"`
	Status  string    `json:"status" validate:"required,oneof=pending paid shipped delivered cancelled"`
}

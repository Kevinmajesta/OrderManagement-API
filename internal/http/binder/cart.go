package binder

import "github.com/google/uuid"

type CartAddItemRequest struct {
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type CartUpdateItemRequest struct {
	Quantity int `json:"quantity"`
}

type CartCheckoutRequest struct {
	UserID        uuid.UUID `json:"user_id"`
	PaymentMethod string    `json:"payment_method"`
	PaidAmount    float64   `json:"paid_amount"`
}

package binder

import "github.com/google/uuid"

type ReceiptGenerateRequest struct {
	OrderID     uuid.UUID `json:"order_id"`
	UserID      uuid.UUID `json:"user_id"`
	CashierName string    `json:"cashier_name"`
}

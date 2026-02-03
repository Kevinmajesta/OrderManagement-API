package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID       uuid.UUID   `json:"order_id" gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID   `json:"user_id" gorm:"column:user_id"`
	TotalPrice    float64     `json:"total_price"`
	PaymentMethod string      `json:"payment_method" gorm:"column:payment_method"`
	PaidAmount    float64     `json:"paid_amount" gorm:"column:paid_amount"`
	ChangeAmount  float64     `json:"change_amount" gorm:"column:change_amount"`
	Status        string      `json:"status" gorm:"default:'pending'"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	OrderItems    []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
	SnapToken     string      `json:"snap_token" gorm:"-"`
	RedirectURL   string      `json:"redirect_url" gorm:"-"`
}

type OrderItem struct {
	OrderItemID  uuid.UUID `json:"order_item_id" gorm:"column:orderitem_id;type:uuid;primaryKey"`
	OrderID      uuid.UUID `json:"order_id"`
	ProductID    uuid.UUID `json:"product_id"`
	Quantity     int       `json:"quantity"`
	PricePerItem float64   `json:"price_per_item"`
	TotalPrice   float64   `json:"total_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

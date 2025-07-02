package entity

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	OrderID    uuid.UUID   `json:"order_id" gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID   `json:"user_id" gorm:"column:user_id"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status" gorm:"default:'pending'"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	OrderItemID uuid.UUID `json:"order_item_id" gorm:"column:orderitem_id;type:uuid;primaryKey"`
	OrderID      uuid.UUID `json:"order_id"`
	ProductID    uuid.UUID `json:"product_id"`
	Quantity     int       `json:"quantity"`
	PricePerItem float64   `json:"price_per_item"`
	TotalPrice   float64   `json:"total_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

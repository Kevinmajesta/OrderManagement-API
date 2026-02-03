package entity

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	CartID    uuid.UUID  `json:"cart_id" gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID  `json:"user_id" gorm:"column:user_id"`
	Status    string     `json:"status" gorm:"default:'active'"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID"`
}

type CartItem struct {
	CartItemID uuid.UUID `json:"cart_item_id" gorm:"type:uuid;primaryKey"`
	CartID     uuid.UUID `json:"cart_id" gorm:"column:cart_id"`
	ProductID  uuid.UUID `json:"product_id" gorm:"column:product_id"`
	Quantity   int       `json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
	ReceiptID     uuid.UUID     `json:"receipt_id" gorm:"type:uuid;primaryKey"`
	OrderID       uuid.UUID     `json:"order_id" gorm:"column:order_id"`
	UserID        uuid.UUID     `json:"user_id" gorm:"column:user_id"`
	Subtotal      float64       `json:"subtotal" gorm:"column:subtotal"`
	TaxAmount     float64       `json:"tax_amount" gorm:"column:tax_amount"`
	TotalAmount   float64       `json:"total_amount" gorm:"column:total_amount"`
	PaymentMethod string        `json:"payment_method" gorm:"column:payment_method"`
	PaymentStatus string        `json:"payment_status" gorm:"column:payment_status"`
	ReceiptNumber string        `json:"receipt_number" gorm:"column:receipt_number;unique"`
	CashierName   string        `json:"cashier_name" gorm:"column:cashier_name"`
	StoreName     string        `json:"store_name" gorm:"column:store_name"`
	StoreAddress  string        `json:"store_address" gorm:"column:store_address"`
	StorePhone    string        `json:"store_phone" gorm:"column:store_phone"`
	ReceiptItems  []ReceiptItem `json:"receipt_items" gorm:"foreignKey:ReceiptID"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type ReceiptItem struct {
	ReceiptItemID uuid.UUID `json:"receipt_item_id" gorm:"type:uuid;primaryKey"`
	ReceiptID     uuid.UUID `json:"receipt_id" gorm:"column:receipt_id"`
	ProductName   string    `json:"product_name" gorm:"column:product_name"`
	Quantity      int       `json:"quantity" gorm:"column:quantity"`
	UnitPrice     float64   `json:"unit_price" gorm:"column:unit_price"`
	TotalPrice    float64   `json:"total_price" gorm:"column:total_price"`
	CreatedAt     time.Time `json:"created_at"`
}

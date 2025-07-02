package repository

import (
	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/pkg/cache"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *entity.Order) error
	UpdateProductStock(productID string, qty int) error
	GetProductByID(productID string) (*entity.Products, error)
	UpdateOrderStatus(orderID uuid.UUID, status string) error
	GetOrderHistoryByUserID(userID string) ([]entity.Order, error)
}

type orderRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewOrderRepository(db *gorm.DB, cacheable cache.Cacheable) OrderRepository {
	return &orderRepository{db: db, cacheable: cacheable}
}

func (r *orderRepository) CreateOrder(order *entity.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) UpdateProductStock(productID string, qty int) error {
	return r.db.Model(&entity.Products{}).
		Where("product_id = ?", productID).
		Update("stock", gorm.Expr("stock - ?", qty)).Error
}

func (r *orderRepository) GetProductByID(productID string) (*entity.Products, error) {
	var product entity.Products
	err := r.db.First(&product, "product_id = ?", productID).Error
	return &product, err
}

func (r *orderRepository) UpdateOrderStatus(orderID uuid.UUID, status string) error {
	return r.db.Model(&entity.Order{}).
		Where("order_id = ?", orderID).
		Update("status", status).Error
}

func (r *orderRepository) GetOrderHistoryByUserID(userID string) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.Preload("OrderItems").Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

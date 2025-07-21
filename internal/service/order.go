package service

import (
	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/internal/repository"
	"errors"

	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(order *entity.Order) error
	UpdateOrderStatus(orderID uuid.UUID, status string) error
	GetOrderHistory(userID string) ([]entity.Order, error)
}

type orderService struct {
	repo repository.OrderRepository // Pastikan ini mengarah ke interface repository Anda
	db   *gorm.DB
}

func NewOrderService(repo repository.OrderRepository, db *gorm.DB) *orderService {
	return &orderService{repo: repo, db: db}
}

// CreateOrder handles the entire order creation process, including stock validation and update.
func (s *orderService) CreateOrder(order *entity.Order) error {
	var totalPrice float64

	// --- Mulai Transaksi di sini ---
	tx := s.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}() // Handle panic during transaction

	// Menggunakan tx untuk operasi database dalam transaksi
	// Pastikan repository Anda menerima *gorm.DB atau *gorm.Tx

	// Loop items, cek stock, hitung harga, dan kurangi stok SECARA ATOMIK
	for i, item := range order.OrderItems {
		var product entity.Products
		// 1. Kunci baris produk dengan FOR UPDATE
		// Pastikan ProductID adalah string UUID yang valid untuk GORM
		err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("product_id = ?", item.ProductID.String()).
			First(&product).Error

		if err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("product not found: " + item.ProductID.String())
			}
			return fmt.Errorf("failed to get product with lock: %w", err)
		}

		// 2. Cek stok setelah baris dikunci
		if product.Stock < item.Quantity {
			tx.Rollback() // Rollback jika stok tidak cukup
			return errors.New("insufficient stock for product: " + product.Name)
		}

		// 3. Update stok di dalam transaksi
		newStock := product.Stock - item.Quantity
		err = tx.Model(&entity.Products{}).
			Where("product_id = ?", item.ProductID.String()).
			Update("stock", newStock).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update product stock: %w", err)
		}

		order.OrderItems[i].PricePerItem = product.Price
		order.OrderItems[i].TotalPrice = float64(item.Quantity) * product.Price
		order.OrderItems[i].OrderItemID = uuid.New() // Generate UUID for OrderItem if needed

		totalPrice += order.OrderItems[i].TotalPrice
	}

	order.OrderID = uuid.New()
	order.TotalPrice = totalPrice
	order.Status = "pending" // Set default status

	// 4. Simpan order dalam transaksi yang sama
	// Anda mungkin perlu memodifikasi `CreateOrder` di repository untuk menerima *gorm.Tx
	err := tx.Create(order).Error // Langsung gunakan tx untuk membuat order
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create order: %w", err)
	}

	// 5. Commit transaksi jika semua berhasil
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *orderService) UpdateOrderStatus(orderID uuid.UUID, status string) error {
	return s.repo.UpdateOrderStatus(orderID, status)
}

func (s *orderService) GetOrderHistory(userID string) ([]entity.Order, error) {
	return s.repo.GetOrderHistoryByUserID(userID)
}

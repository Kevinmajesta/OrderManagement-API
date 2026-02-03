package service

import (
	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/internal/repository"
	"Kevinmajesta/OrderManagementAPI/pkg/midtrans"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(order *entity.Order) error
	UpdateOrderStatus(orderID uuid.UUID, status string) error
	UpdateOrderStatusByOrderID(orderID string, status string) error
	GetOrderHistory(userID string) ([]entity.Order, error)
}

type orderService struct {
	repo            repository.OrderRepository
	db              *gorm.DB
	midtransService *midtrans.MidtransService
}

func NewOrderService(repo repository.OrderRepository, db *gorm.DB, midtransService *midtrans.MidtransService) *orderService {
	return &orderService{
		repo:            repo,
		db:              db,
		midtransService: midtransService,
	}
}

func (s *orderService) CreateOrder(order *entity.Order) error {
	var totalPrice float64

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i, item := range order.OrderItems {
		// Generate UUID untuk OrderItem
		order.OrderItems[i].OrderItemID = uuid.New()
		
		var product entity.Products
		err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("product_id = ?", item.ProductID.String()).
			First(&product).Error

		if err != nil {
			tx.Rollback()
			return err
		}

		if product.Stock < item.Quantity {
			tx.Rollback()
			return errors.New("stok tidak cukup")
		}

		tx.Model(&product).Update("stock", product.Stock-item.Quantity)

		order.OrderItems[i].PricePerItem = product.Price
		order.OrderItems[i].TotalPrice = float64(item.Quantity) * product.Price
		totalPrice += order.OrderItems[i].TotalPrice
	}

	order.OrderID = uuid.New()
	order.TotalPrice = totalPrice
	order.Status = "pending"

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction first
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// Get user details for Midtrans
	var user entity.User
	if err := s.db.Where("user_id = ?", order.UserID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to get user details: %v", err)
	}

	// Create Midtrans transaction (after order saved)
	snapResp, errMidtrans := s.midtransService.CreateTransaction(
		order.OrderID.String(),
		int64(totalPrice),
		user.Fullname,
		user.Email,
		user.Phone,
	)
	if errMidtrans != nil {
		return fmt.Errorf("midtrans error: %v", errMidtrans)
	}

	// Set to order object for response (tidak disimpan ke DB)
	order.SnapToken = snapResp.Token
	order.RedirectURL = snapResp.RedirectURL

	return nil
}

func (s *orderService) UpdateOrderStatus(orderID uuid.UUID, status string) error {
	return s.repo.UpdateOrderStatus(orderID, status)
}

func (s *orderService) UpdateOrderStatusByOrderID(orderID string, status string) error {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return fmt.Errorf("invalid order_id format: %v", err)
	}
	return s.repo.UpdateOrderStatus(orderUUID, status)
}

func (s *orderService) GetOrderHistory(userID string) ([]entity.Order, error) {
	return s.repo.GetOrderHistoryByUserID(userID)
}

package service

import (
	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/internal/repository"
	"errors"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(order *entity.Order) error
	UpdateOrderStatus(orderID uuid.UUID, status string) error
	GetOrderHistory(userID string) ([]entity.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(order *entity.Order) error {
	var totalPrice float64

	// Loop items, cek stock, hitung harga
	for i, item := range order.OrderItems {
		product, err := s.repo.GetProductByID(item.ProductID.String())
		if err != nil {
			return errors.New("product not found")
		}
		if product.Stock < item.Quantity {
			return errors.New("insufficient stock for product: " + product.Name)
		}

		order.OrderItems[i].PricePerItem = product.Price
		order.OrderItems[i].TotalPrice = float64(item.Quantity) * product.Price
		order.OrderItems[i].OrderItemID = uuid.New()

		totalPrice += order.OrderItems[i].TotalPrice
	}

	order.OrderID = uuid.New()
	order.TotalPrice = totalPrice

	// Simpan order dan update stok (dalam transaksi)
	err := s.repo.CreateOrder(order)
	if err != nil {
		return err
	}

	// Kurangi stok
	for _, item := range order.OrderItems {
		err := s.repo.UpdateProductStock(item.ProductID.String(), item.Quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *orderService) UpdateOrderStatus(orderID uuid.UUID, status string) error {
	return s.repo.UpdateOrderStatus(orderID, status)
}

func (s *orderService) GetOrderHistory(userID string) ([]entity.Order, error) {
	return s.repo.GetOrderHistoryByUserID(userID)
}

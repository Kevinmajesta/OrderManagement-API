package service

import (
	"errors"
	"fmt"
	"time"

	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReceiptService interface {
	GenerateReceipt(orderID uuid.UUID, userID uuid.UUID, cashierName string) (*entity.Receipt, error)
	GetReceiptByID(receiptID uuid.UUID) (*entity.Receipt, error)
	GetReceiptByOrderID(orderID uuid.UUID) (*entity.Receipt, error)
	GetReceiptsByUserID(userID uuid.UUID) ([]entity.Receipt, error)
}

type receiptService struct {
	receiptRepo repository.ReceiptRepository
	orderRepo   repository.OrderRepository
	db          *gorm.DB
}

func NewReceiptService(receiptRepo repository.ReceiptRepository, orderRepo repository.OrderRepository, db *gorm.DB) *receiptService {
	return &receiptService{
		receiptRepo: receiptRepo,
		orderRepo:   orderRepo,
		db:          db,
	}
}

func (s *receiptService) GenerateReceipt(orderID uuid.UUID, userID uuid.UUID, cashierName string) (*entity.Receipt, error) {
	var order entity.Order
	err := s.db.Preload("OrderItems").Where("order_id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Check if receipt already exists
	existing, _ := s.receiptRepo.GetReceiptByOrderID(orderID)
	if existing != nil {
		return existing, nil
	}

	// Generate receipt number
	lastNum, _ := s.receiptRepo.GetLastReceiptNumber()
	receiptNumber := s.generateReceiptNumber(lastNum)

	// Calculate tax (10%)
	taxRate := 0.10
	taxAmount := order.TotalPrice * taxRate
	totalAmount := order.TotalPrice + taxAmount

	receipt := &entity.Receipt{
		ReceiptID:     uuid.New(),
		OrderID:       orderID,
		UserID:        userID,
		Subtotal:      order.TotalPrice,
		TaxAmount:     taxAmount,
		TotalAmount:   totalAmount,
		PaymentMethod: order.PaymentMethod,
		PaymentStatus: order.Status,
		ReceiptNumber: receiptNumber,
		CashierName:   cashierName,
		StoreName:     "Cuaniaga Store",
		StoreAddress:  "Jl. Raya No. 123",
		StorePhone:    "+62 812 3456 7890",
		ReceiptItems:  make([]entity.ReceiptItem, 0),
	}

	// Add receipt items
	for _, item := range order.OrderItems {
		receiptItem := entity.ReceiptItem{
			ReceiptItemID: uuid.New(),
			ReceiptID:     receipt.ReceiptID,
			ProductName:   s.getProductName(item.ProductID),
			Quantity:      item.Quantity,
			UnitPrice:     item.PricePerItem,
			TotalPrice:    item.TotalPrice,
		}
		receipt.ReceiptItems = append(receipt.ReceiptItems, receiptItem)
	}

	if err := s.receiptRepo.CreateReceipt(receipt); err != nil {
		return nil, err
	}

	return receipt, nil
}

func (s *receiptService) GetReceiptByID(receiptID uuid.UUID) (*entity.Receipt, error) {
	return s.receiptRepo.GetReceiptByID(receiptID)
}

func (s *receiptService) GetReceiptByOrderID(orderID uuid.UUID) (*entity.Receipt, error) {
	return s.receiptRepo.GetReceiptByOrderID(orderID)
}

func (s *receiptService) GetReceiptsByUserID(userID uuid.UUID) ([]entity.Receipt, error) {
	return s.receiptRepo.GetReceiptsByUserID(userID)
}

func (s *receiptService) generateReceiptNumber(lastNum string) string {
	now := time.Now()
	date := now.Format("20060102")

	if lastNum == "" {
		return fmt.Sprintf("RCP%s0001", date)
	}

	lastDate := lastNum[3:11]
	if lastDate != date {
		return fmt.Sprintf("RCP%s0001", date)
	}

	var counter int
	fmt.Sscanf(lastNum, fmt.Sprintf("RCP%s%%04d", date), &counter)
	counter++
	return fmt.Sprintf("RCP%s%04d", date, counter)
}

func (s *receiptService) getProductName(productID uuid.UUID) string {
	product, err := s.orderRepo.GetProductByID(productID.String())
	if err != nil {
		return "Unknown Product"
	}
	return product.Name
}

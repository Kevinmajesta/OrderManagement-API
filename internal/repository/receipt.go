package repository

import (
	"errors"

	"Kevinmajesta/OrderManagementAPI/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReceiptRepository interface {
	CreateReceipt(receipt *entity.Receipt) error
	GetReceiptByID(receiptID uuid.UUID) (*entity.Receipt, error)
	GetReceiptByOrderID(orderID uuid.UUID) (*entity.Receipt, error)
	GetReceiptsByUserID(userID uuid.UUID) ([]entity.Receipt, error)
	GetLastReceiptNumber() (string, error)
}

type receiptRepository struct {
	db *gorm.DB
}

func NewReceiptRepository(db *gorm.DB) ReceiptRepository {
	return &receiptRepository{db: db}
}

func (r *receiptRepository) CreateReceipt(receipt *entity.Receipt) error {
	if receipt == nil {
		return errors.New("receipt is nil")
	}
	return r.db.Create(receipt).Error
}

func (r *receiptRepository) GetReceiptByID(receiptID uuid.UUID) (*entity.Receipt, error) {
	var receipt entity.Receipt
	err := r.db.Preload("ReceiptItems").Where("receipt_id = ?", receiptID).First(&receipt).Error
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}

func (r *receiptRepository) GetReceiptByOrderID(orderID uuid.UUID) (*entity.Receipt, error) {
	var receipt entity.Receipt
	err := r.db.Preload("ReceiptItems").Where("order_id = ?", orderID).First(&receipt).Error
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}

func (r *receiptRepository) GetReceiptsByUserID(userID uuid.UUID) ([]entity.Receipt, error) {
	var receipts []entity.Receipt
	err := r.db.Preload("ReceiptItems").Where("user_id = ?", userID).Order("created_at DESC").Find(&receipts).Error
	if err != nil {
		return nil, err
	}
	return receipts, nil
}

func (r *receiptRepository) GetLastReceiptNumber() (string, error) {
	var receipt entity.Receipt
	err := r.db.Order("created_at DESC").First(&receipt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return receipt.ReceiptNumber, nil
}

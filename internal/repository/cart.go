package repository

import (
	"errors"

	"Kevinmajesta/OrderManagementAPI/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetActiveCartByUserID(userID uuid.UUID) (*entity.Cart, error)
	CreateCart(userID uuid.UUID) (*entity.Cart, error)
	GetCartWithItems(cartID uuid.UUID) (*entity.Cart, error)
	GetCartItem(cartID uuid.UUID, productID uuid.UUID) (*entity.CartItem, error)
	GetCartItemByID(cartItemID uuid.UUID) (*entity.CartItem, error)
	CreateCartItem(item *entity.CartItem) error
	UpdateCartItemQuantity(cartItemID uuid.UUID, qty int) error
	DeleteCartItem(cartItemID uuid.UUID) error
	SetCartStatus(cartID uuid.UUID, status string) error
	ClearCartItems(cartID uuid.UUID) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetActiveCartByUserID(userID uuid.UUID) (*entity.Cart, error) {
	var cart entity.Cart
	err := r.db.Preload("Items").Where("user_id = ? AND status = ?", userID, "active").First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) CreateCart(userID uuid.UUID) (*entity.Cart, error) {
	cart := &entity.Cart{
		CartID: uuid.New(),
		UserID: userID,
		Status: "active",
	}
	if err := r.db.Create(cart).Error; err != nil {
		return nil, err
	}
	return cart, nil
}

func (r *cartRepository) GetCartWithItems(cartID uuid.UUID) (*entity.Cart, error) {
	var cart entity.Cart
	err := r.db.Preload("Items").Where("cart_id = ?", cartID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetCartItem(cartID uuid.UUID, productID uuid.UUID) (*entity.CartItem, error) {
	var item entity.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) GetCartItemByID(cartItemID uuid.UUID) (*entity.CartItem, error) {
	var item entity.CartItem
	err := r.db.Where("cart_item_id = ?", cartItemID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) CreateCartItem(item *entity.CartItem) error {
	if item == nil {
		return errors.New("cart item is nil")
	}
	return r.db.Create(item).Error
}

func (r *cartRepository) UpdateCartItemQuantity(cartItemID uuid.UUID, qty int) error {
	return r.db.Model(&entity.CartItem{}).Where("cart_item_id = ?", cartItemID).Update("quantity", qty).Error
}

func (r *cartRepository) DeleteCartItem(cartItemID uuid.UUID) error {
	return r.db.Where("cart_item_id = ?", cartItemID).Delete(&entity.CartItem{}).Error
}

func (r *cartRepository) SetCartStatus(cartID uuid.UUID, status string) error {
	return r.db.Model(&entity.Cart{}).Where("cart_id = ?", cartID).Update("status", status).Error
}

func (r *cartRepository) ClearCartItems(cartID uuid.UUID) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&entity.CartItem{}).Error
}

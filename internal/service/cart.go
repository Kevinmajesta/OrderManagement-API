package service

import (
	"errors"

	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartService interface {
	AddItem(userID uuid.UUID, productID uuid.UUID, qty int) (*entity.Cart, error)
	UpdateItem(cartItemID uuid.UUID, qty int) (*entity.Cart, error)
	RemoveItem(cartItemID uuid.UUID) error
	GetCart(userID uuid.UUID) (*entity.Cart, error)
	Checkout(userID uuid.UUID, paymentMethod string, paidAmount float64) (*entity.Order, error)
}

type cartService struct {
	cartRepository repository.CartRepository
	orderService   OrderService
	productRepo    repository.ProductRepository
}

func NewCartService(cartRepository repository.CartRepository, orderService OrderService, productRepo repository.ProductRepository) *cartService {
	return &cartService{
		cartRepository: cartRepository,
		orderService:   orderService,
		productRepo:    productRepo,
	}
}

func (s *cartService) AddItem(userID uuid.UUID, productID uuid.UUID, qty int) (*entity.Cart, error) {
	if qty <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	// Ensure product exists
	if _, err := s.productRepo.FindProductByID(productID.String()); err != nil {
		return nil, errors.New("product not found")
	}

	cart, err := s.cartRepository.GetActiveCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart, err = s.cartRepository.CreateCart(userID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	item, err := s.cartRepository.GetCartItem(cart.CartID, productID)
	if err == nil && item != nil {
		newQty := item.Quantity + qty
		if err := s.cartRepository.UpdateCartItemQuantity(item.CartItemID, newQty); err != nil {
			return nil, err
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		newItem := &entity.CartItem{
			CartItemID: uuid.New(),
			CartID:     cart.CartID,
			ProductID:  productID,
			Quantity:   qty,
		}
		if err := s.cartRepository.CreateCartItem(newItem); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return s.cartRepository.GetCartWithItems(cart.CartID)
}

func (s *cartService) UpdateItem(cartItemID uuid.UUID, qty int) (*entity.Cart, error) {
	if qty <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	if err := s.cartRepository.UpdateCartItemQuantity(cartItemID, qty); err != nil {
		return nil, err
	}

	item, err := s.cartRepository.GetCartItemByID(cartItemID)
	if err != nil {
		return nil, err
	}

	return s.cartRepository.GetCartWithItems(item.CartID)
}

func (s *cartService) RemoveItem(cartItemID uuid.UUID) error {
	return s.cartRepository.DeleteCartItem(cartItemID)
}

func (s *cartService) GetCart(userID uuid.UUID) (*entity.Cart, error) {
	cart, err := s.cartRepository.GetActiveCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.cartRepository.CreateCart(userID)
		}
		return nil, err
	}
	return cart, nil
}

func (s *cartService) Checkout(userID uuid.UUID, paymentMethod string, paidAmount float64) (*entity.Order, error) {
	cart, err := s.cartRepository.GetActiveCartByUserID(userID)
	if err != nil {
		return nil, err
	}
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	orderItems := make([]entity.OrderItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		orderItems = append(orderItems, entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	order := &entity.Order{
		UserID:        userID,
		PaymentMethod: paymentMethod,
		PaidAmount:    paidAmount,
		OrderItems:    orderItems,
	}

	if err := s.orderService.CreateOrder(order); err != nil {
		return nil, err
	}

	if err := s.cartRepository.SetCartStatus(cart.CartID, "checked_out"); err != nil {
		return nil, err
	}

	return order, nil
}

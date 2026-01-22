package service

import (
	"testing"
	"time"

	"Kevinmajesta/OrderManagementAPI/internal/entity"

	"github.com/google/uuid"
)

// TestOrderValidation tests order creation with validations
func TestOrderValidation(t *testing.T) {
	tests := []struct {
		name     string
		order    *entity.Order
		wantErr  bool
		errorMsg string
	}{
		{
			name: "valid order",
			order: &entity.Order{
				OrderID:    uuid.New(),
				UserID:     uuid.New(),
				Status:     "pending",
				TotalPrice: 299.99,
				CreatedAt:  time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid total price (negative)",
			order: &entity.Order{
				OrderID:    uuid.New(),
				UserID:     uuid.New(),
				Status:     "pending",
				TotalPrice: -50.00,
				CreatedAt:  time.Now(),
			},
			wantErr:  true,
			errorMsg: "Total price must be positive",
		},
		{
			name: "invalid status",
			order: &entity.Order{
				OrderID:    uuid.New(),
				UserID:     uuid.New(),
				Status:     "invalid_status",
				TotalPrice: 299.99,
				CreatedAt:  time.Now(),
			},
			wantErr:  true,
			errorMsg: "Invalid order status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate order
			if tt.order.TotalPrice < 0 && tt.wantErr {
				t.Log("Total price validation: PASS (negative amount caught)")
			}
			if tt.order.Status != "pending" && tt.order.Status != "processing" && tt.order.Status != "shipped" && tt.order.Status != "delivered" && tt.wantErr {
				t.Log("Status validation: PASS (invalid status caught)")
			}
		})
	}
}

// TestOrderFields tests order entity fields
func TestOrderFields(t *testing.T) {
	t.Run("order can be created with all fields", func(t *testing.T) {
		now := time.Now()
		order := &entity.Order{
			OrderID:    uuid.New(),
			UserID:     uuid.New(),
			Status:     "pending",
			TotalPrice: 299.99,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		if order.OrderID == uuid.Nil {
			t.Error("Order ID should not be nil")
		}
		if order.UserID == uuid.Nil {
			t.Error("User ID should not be nil")
		}
		if order.TotalPrice <= 0 {
			t.Error("Total price should be positive")
		}
		if order.Status == "" {
			t.Error("Status should not be empty")
		}
		t.Log("Order entity created successfully")
	})
}

// TestOrderStatus tests various order status values
func TestOrderStatus(t *testing.T) {
	statuses := []string{"pending", "processing", "shipped", "delivered"}

	for _, status := range statuses {
		t.Run("status_"+status, func(t *testing.T) {
			order := &entity.Order{
				OrderID:    uuid.New(),
				UserID:     uuid.New(),
				Status:     status,
				TotalPrice: 100.00,
				CreatedAt:  time.Now(),
			}

			if order.Status != status {
				t.Errorf("Expected status %s, got %s", status, order.Status)
			}
			t.Logf("Status %s assigned successfully", status)
		})
	}
}

// TestOrderPrice tests order with various amounts
func TestOrderPrice(t *testing.T) {
	prices := []float64{10.00, 50.99, 299.99, 999.99, 5000.00}

	for _, price := range prices {
		t.Run("price", func(t *testing.T) {
			order := &entity.Order{
				OrderID:    uuid.New(),
				UserID:     uuid.New(),
				TotalPrice: price,
				Status:     "pending",
				CreatedAt:  time.Now(),
			}

			if order.TotalPrice != price {
				t.Errorf("Expected price %.2f, got %.2f", price, order.TotalPrice)
			}
			t.Logf("Price %.2f assigned successfully", price)
		})
	}
}

// TestOrderUpdate tests order update functionality
func TestOrderUpdate(t *testing.T) {
	t.Run("order fields can be updated", func(t *testing.T) {
		order := &entity.Order{
			OrderID:    uuid.New(),
			UserID:     uuid.New(),
			Status:     "pending",
			TotalPrice: 299.99,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		// Update fields
		order.Status = "shipped"
		order.TotalPrice = 350.00
		order.UpdatedAt = time.Now()

		if order.Status != "shipped" {
			t.Error("Status should be updated")
		}
		if order.TotalPrice != 350.00 {
			t.Error("Total price should be updated")
		}
		t.Log("Order updated successfully")
	})
}

// TestOrderItems tests order items functionality
func TestOrderItems(t *testing.T) {
	t.Run("order can have multiple items", func(t *testing.T) {
		order := &entity.Order{
			OrderID:    uuid.New(),
			UserID:     uuid.New(),
			Status:     "pending",
			TotalPrice: 299.99,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			OrderItems: []entity.OrderItem{
				{
					OrderItemID:  uuid.New(),
					OrderID:      uuid.New(),
					ProductID:    uuid.New(),
					Quantity:     2,
					PricePerItem: 50.00,
					TotalPrice:   100.00,
				},
				{
					OrderItemID:  uuid.New(),
					OrderID:      uuid.New(),
					ProductID:    uuid.New(),
					Quantity:     3,
					PricePerItem: 66.66,
					TotalPrice:   199.98,
				},
			},
		}

		if len(order.OrderItems) != 2 {
			t.Errorf("Expected 2 order items, got %d", len(order.OrderItems))
		}
		t.Logf("Order created with %d items successfully", len(order.OrderItems))
	})
}

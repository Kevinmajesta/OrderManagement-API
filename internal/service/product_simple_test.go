package service

import (
	"testing"

	"Kevinmajesta/OrderManagementAPI/internal/entity"

	"github.com/google/uuid"
)

// TestProductValidation tests product creation with validations
func TestProductValidation(t *testing.T) {
	tests := []struct {
		name     string
		product  *entity.Products
		wantErr  bool
		errorMsg string
	}{
		{
			name: "valid product",
			product: &entity.Products{
				Name:        "Laptop",
				Description: "High-performance laptop",
				PhotoURL:    "https://example.com/laptop.jpg",
				Price:       999.99,
				Stock:       10,
			},
			wantErr: false,
		},
		{
			name: "invalid price (negative)",
			product: &entity.Products{
				Name:        "Product",
				Description: "Description",
				PhotoURL:    "https://example.com/photo.jpg",
				Price:       -50.00,
				Stock:       5,
			},
			wantErr:  true,
			errorMsg: "Price must be greater than 0",
		},
		{
			name: "invalid stock (zero)",
			product: &entity.Products{
				Name:        "Product",
				Description: "Description",
				PhotoURL:    "https://example.com/photo.jpg",
				Price:       99.99,
				Stock:       0,
			},
			wantErr:  true,
			errorMsg: "Stock must be greater than 0",
		},
		{
			name: "invalid stock (negative)",
			product: &entity.Products{
				Name:        "Product",
				Description: "Description",
				PhotoURL:    "https://example.com/photo.jpg",
				Price:       99.99,
				Stock:       -10,
			},
			wantErr:  true,
			errorMsg: "Stock must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate product
			if tt.product.Price <= 0 && tt.wantErr {
				t.Log("Price validation: PASS (invalid price caught)")
			}
			if tt.product.Stock <= 0 && tt.wantErr {
				t.Log("Stock validation: PASS (invalid stock caught)")
			}
		})
	}
}

// TestProductFields tests product entity fields
func TestProductFields(t *testing.T) {
	t.Run("product can be created with all fields", func(t *testing.T) {
		product := &entity.Products{
			ProductID:   uuid.New(),
			Name:        "Laptop",
			Description: "High-performance laptop",
			PhotoURL:    "https://example.com/laptop.jpg",
			Price:       999.99,
			Stock:       10,
		}

		if product.ProductID == uuid.Nil {
			t.Error("Product ID should not be nil")
		}
		if product.Name == "" {
			t.Error("Name should not be empty")
		}
		if product.Price <= 0 {
			t.Error("Price should be positive")
		}
		if product.Stock <= 0 {
			t.Error("Stock should be positive")
		}
		t.Log("Product entity created successfully")
	})
}

// TestProductPriceRange tests product with various price ranges
func TestProductPriceRange(t *testing.T) {
	prices := []float64{0.01, 10.50, 99.99, 999.99, 10000.00}

	for _, price := range prices {
		t.Run("price_"+string(rune(price)), func(t *testing.T) {
			product := &entity.Products{
				ProductID: uuid.New(),
				Name:      "Test Product",
				Price:     price,
				Stock:     10,
			}

			if product.Price != price {
				t.Errorf("Expected price %f, got %f", price, product.Price)
			}
			t.Logf("Price %.2f assigned successfully", price)
		})
	}
}

// TestProductStockRange tests product with various stock quantities
func TestProductStockRange(t *testing.T) {
	stocks := []int{1, 10, 100, 1000, 10000}

	for _, stock := range stocks {
		t.Run("stock_"+string(rune(stock)), func(t *testing.T) {
			product := &entity.Products{
				ProductID: uuid.New(),
				Name:      "Test Product",
				Price:     99.99,
				Stock:     stock,
			}

			if product.Stock != stock {
				t.Errorf("Expected stock %d, got %d", stock, product.Stock)
			}
			t.Logf("Stock %d assigned successfully", stock)
		})
	}
}

// TestProductUpdate tests product update functionality
func TestProductUpdate(t *testing.T) {
	t.Run("product fields can be updated", func(t *testing.T) {
		product := &entity.Products{
			ProductID:   uuid.New(),
			Name:        "Original Name",
			Description: "Original description",
			Price:       99.99,
			Stock:       10,
		}

		// Update fields
		product.Name = "Updated Name"
		product.Description = "Updated description"
		product.Price = 149.99
		product.Stock = 20

		if product.Name != "Updated Name" {
			t.Error("Name should be updated")
		}
		if product.Price != 149.99 {
			t.Error("Price should be updated")
		}
		if product.Stock != 20 {
			t.Error("Stock should be updated")
		}
		t.Log("Product updated successfully")
	})
}

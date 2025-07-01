package repository

import (
	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/pkg/cache"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *entity.Products) (*entity.Products, error)
	UpdateProduct(product *entity.Products) (*entity.Products, error)
	CheckProductExists(productId string) (bool, error)
	FindProductByID(productId string) (*entity.Products, error)
	DeleteProduct(product *entity.Products) (bool, error)
	FindAllProduct(page int, search string) ([]entity.Products, error)
}

type productRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewProductRepository(db *gorm.DB, cacheable cache.Cacheable) *productRepository {
	return &productRepository{db: db, cacheable: cacheable}
}

func (r *productRepository) CreateProduct(product *entity.Products) (*entity.Products, error) {
	if err := r.db.Create(&product).Error; err != nil {
		return product, err
	}
	r.cacheable.Delete("FindAllProducts_page_1")
	r.cacheable.Delete("FindAllProducts_page_2")
	return product, nil
}

func (r *productRepository) UpdateProduct(product *entity.Products) (*entity.Products, error) {
	fields := make(map[string]interface{})

	if product.Name != "" {
		fields["name"] = product.Name
	}
	if product.Price != 0 {
		fields["price"] = product.Price
	}
	if product.PhotoURL != "" {
		fields["photo_url"] = product.PhotoURL
	}
	if product.Description != "" {
		fields["description"] = product.Description
	}
	if product.Stock != 0 {
		fields["stock"] = product.Stock
	}

	if err := r.db.Model(product).Where("product_id = ?", product.ProductID).Updates(fields).Error; err != nil {
		return product, err
	}
	r.cacheable.Delete("FindAllProducts_page_1")

	return product, nil
}

func (r *productRepository) CheckProductExists(productId string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Products{}).Where("product_id = ?", productId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *productRepository) FindProductByID(productId string) (*entity.Products, error) {
	product := new(entity.Products)
	if err := r.db.Where("product_id = ?", productId).First(product).Error; err != nil {
		log.Printf("Error finding product by ID: %v", err)
		return nil, err // Pastikan mengembalikan nil, err
	}
	log.Printf("Product found: %v", product)
	return product, nil
}

func (r *productRepository) DeleteProduct(product *entity.Products) (bool, error) {
	log.Printf("Deleting product: %v", product)
	// Ensure hard delete by using Unscoped()
	if err := r.db.Unscoped().Delete(product).Error; err != nil {
		log.Printf("Error deleting product: %v", err)
		return false, err
	}
	log.Println("Product deleted successfully")
	r.cacheable.Delete("FindAllProducts_page_1")
	return true, nil
}

func (r *productRepository) FindAllProduct(page int, search string) ([]entity.Products, error) {
	var products []entity.Products
	const pageSize = 100
	offset := (page - 1) * pageSize

	// Jika ada keyword pencarian, jangan pakai cache
	if search != "" {
		query := r.db.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", "%"+strings.ToLower(search)+"%", "%"+strings.ToLower(search)+"%")
		if err := query.Limit(pageSize).Offset(offset).Find(&products).Error; err != nil {
			return products, err
		}
		return products, nil
	}

	// Tanpa pencarian → pakai cache
	key := fmt.Sprintf("FindAllProducts_page_%d", page)
	data, _ := r.cacheable.Get(key)
	if data == "" {
		if err := r.db.Limit(pageSize).Offset(offset).Find(&products).Error; err != nil {
			return products, err
		}
		marshalled, _ := json.Marshal(products)
		_ = r.cacheable.Set(key, marshalled, 5*time.Minute)
	} else {
		_ = json.Unmarshal([]byte(data), &products)
	}

	return products, nil
}

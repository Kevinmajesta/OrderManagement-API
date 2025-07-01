package repository

import (
	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/pkg/cache"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *entity.Products) (*entity.Products, error)
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

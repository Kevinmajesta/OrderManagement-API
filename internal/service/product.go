package service

import (
	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/internal/repository"
	"errors"
)

type ProductService interface {
	CreateProduct(product *entity.Products) (*entity.Products, error)
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) *productService {
	return &productService{
		productRepository: productRepository,
	}
}

func (s *productService) CreateProduct(product *entity.Products) (*entity.Products, error) {
	if product.Name == "" {
		return nil, errors.New("Product name cannot be empty")
	}

	if product.Description == "" {
		return nil, errors.New("Product description cannot be empty")
	}

	if product.PhotoURL == "" {
		return nil, errors.New("Photo URL cannot be empty")
	}

	if product.Price <= 0 {
		return nil, errors.New("Price must be greater than 0")
	}

	if product.Stock <= 0 {
		return nil, errors.New("Stock must be greater than 0")
	}

	newProduct := entity.NewProduct(
		product.Name,
		product.Description,
		product.PhotoURL,
		product.Price,
		product.Stock,
	)

	savedProduct, err := s.productRepository.CreateProduct(newProduct)
	if err != nil {
		return nil, err
	}

	return savedProduct, nil
}

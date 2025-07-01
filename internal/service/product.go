package service

import (
	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/internal/repository"
	"errors"
	"log"
)

type ProductService interface {
	CreateProduct(product *entity.Products) (*entity.Products, error)
	UpdateProduct(product *entity.Products) (*entity.Products, error)
	CheckProductExists(productId string) (bool, error)
	FindProductByID(productId string) (*entity.Products, error)
	DeleteProduct(productId string) (bool, error)
	FindAllProduct(page int, search string) ([]entity.Products, error)
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

func (s *productService) UpdateProduct(product *entity.Products) (*entity.Products, error) {
	if product.Name == "" {
		return nil, errors.New("Product name cannot be empty")
	}

	if product.Description == "" {
		return nil, errors.New("Product description cannot be empty")
	}

	if product.Price <= 0 {
		return nil, errors.New("Price must be greater than 0")
	}

	if product.Stock <= 0 {
		return nil, errors.New("Stock must be greater than 0")
	}

	updatedProduct, err := s.productRepository.UpdateProduct(product)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (s *productService) CheckProductExists(productId string) (bool, error) {
	return s.productRepository.CheckProductExists(productId)
}

func (s *productService) FindProductByID(productId string) (*entity.Products, error) {
	return s.productRepository.FindProductByID(productId)
}

func (s *productService) DeleteProduct(productId string) (bool, error) {
	product, err := s.productRepository.FindProductByID(productId)
	if err != nil {
		return false, err
	}

	log.Printf("Product to be deleted: %v", product)
	return s.productRepository.DeleteProduct(product)
}

func (s *productService) FindAllProduct(page int, search string) ([]entity.Products, error) {
	return s.productRepository.FindAllProduct(page, search)
}

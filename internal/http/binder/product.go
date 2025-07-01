package binder

import (
	"mime/multipart"
	"github.com/google/uuid"
)

type ProductCreateRequest struct {
	Name        string                `form:"name" json:"name" validate:"required"`
	Description string                `form:"description" json:"description"`
	Photo       *multipart.FileHeader `form:"photo" json:"-" validate:"required"`
	Price       float64               `form:"price" json:"price" validate:"required,min=0"`
	Stock       int                   `form:"stock" json:"stock" validate:"required,min=0"`
}

type ProductUpdateRequest struct {
	ProductID   uuid.UUID             `param:"id" json:"id" validate:"required"`
	Name        string                `form:"name" json:"name" validate:"required"`
	Description string                `form:"description" json:"description"`
	Photo       *multipart.FileHeader `form:"photo" json:"-"`
	Price       float64               `form:"price" json:"price" validate:"required,min=0"`
	Stock       int                   `form:"stock" json:"stock" validate:"required,min=0"`
}

type ProductDeleteRequest struct {
	ProductID uuid.UUID `param:"id" validate:"required"`
}

type ProductUpdateStockRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required"`
}


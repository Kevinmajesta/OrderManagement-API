package entity

import (
	"github.com/google/uuid"
)

type Products struct {
	ProductID   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"column:name;not null;unique"`
	Description string    `json:"description" gorm:"column:description"`
	PhotoURL    string    `json:"photo_url" gorm:"column:photo_url"`
	Price       float64   `json:"price" gorm:"column:price;type:numeric(10,2);not null;check:price >= 0"`
	Stock       int       `json:"stock" gorm:"column:stock;type:integer;not null;default:0;check:stock >= 0"`
	Auditable
}

func NewProduct(name, description, photoURL string, price float64, stock int) *Products {
	return &Products{
		Name:        name,
		Description: description,
		PhotoURL:    photoURL,
		Price:       price,
		Stock:       stock,
		Auditable:   NewAuditable(),
	}
}

func UpdateProduct(productID uuid.UUID, name, description, photoURL string, price float64, stock int) *Products {
	return &Products{
		ProductID:   productID,
		Name:        name,
		Description: description,
		PhotoURL:    photoURL,
		Price:       price,
		Stock:       stock,
		Auditable:   UpdateAuditable(),
	}
}

package product

import "github.com/shohann/golang-ecommerce-api/domain"

type Service interface {
	Create(product domain.Product) (*domain.Product, error)
	GetPublicProducts(page, limit int64) ([]domain.Product, int64, error)
	Delete(id int64) error
}

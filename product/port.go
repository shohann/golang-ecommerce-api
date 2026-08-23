package product

import "github.com/shohann/golang-ecommerce-api/domain"

type ProductRepo interface {
	Create(product domain.Product) (*domain.Product, error)
	List(limit, offset int64) ([]domain.Product, error)
	Count() (int64, error)
	Delete(id int64) error
}

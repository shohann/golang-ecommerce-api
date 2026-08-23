package category

import "github.com/shohann/golang-ecommerce-api/domain"

type Service interface {
	Create(name string) (*domain.Category, error)
	Get(id int64) (*domain.Category, error)
	Update(id int64, name string) (*domain.Category, error)
	Delete(id int64) error
	List(page, limit int64) ([]domain.Category, int64, error)
}

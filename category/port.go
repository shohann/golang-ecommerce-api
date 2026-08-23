package category

import "github.com/shohann/golang-ecommerce-api/domain"

type CategoryRepo interface {
	Create(category domain.Category) (*domain.Category, error)
	FindByID(id int64) (*domain.Category, error)
	Update(category domain.Category) (*domain.Category, error)
	Delete(id int64) error
	List(limit, offset int64) ([]domain.Category, error)
	Count() (int64, error)
	ExistsByName(name string, excludeID int64) (bool, error)
}

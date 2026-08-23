package product

import (
	"github.com/shohann/golang-ecommerce-api/apperr"
	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/domain"
	productHandler "github.com/shohann/golang-ecommerce-api/rest/handlers/product"
)

type service struct {
	productRepo ProductRepo
	cnf         *config.Config
}

type Service interface {
	productHandler.Service
}

func NewService(productRepo ProductRepo, cnf *config.Config) Service {
	return &service{
		productRepo: productRepo,
		cnf:         cnf,
	}
}

func (svc *service) Create(product domain.Product) (*domain.Product, error) {
	createdProduct, err := svc.productRepo.Create(product)
	if err != nil {
		return nil, apperr.WrapInternal("create_product", err)
	}

	if createdProduct == nil {
		return nil, nil
	}

	return &product, nil
}

func (svc *service) GetPublicProducts(page, limit int64) ([]domain.Product, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	total, err := svc.productRepo.Count()

	if err != nil {
		return nil, 0, apperr.WrapInternal("count_products", err)
	}

	products, err := svc.productRepo.List(limit, offset)
	if err != nil {
		return nil, 0, apperr.WrapInternal("list users", err)
	}

	return products, total, nil
}

func (svc *service) Delete(id int64) error {
	if id <= 0 {
		return apperr.Validation("invalid product id")
	}

	err := svc.productRepo.Delete(id)

	if err != nil {
		return apperr.WrapInternal("delete_product", err)
	}

	return nil
}

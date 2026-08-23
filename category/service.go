package category

import (
	"strings"

	"github.com/shohann/golang-ecommerce-api/apperr"
	"github.com/shohann/golang-ecommerce-api/domain"
	categoryHandler "github.com/shohann/golang-ecommerce-api/rest/handlers/category"
)

type service struct {
	categoryRepo CategoryRepo
}

type Service interface {
	categoryHandler.Service
}

func NewService(categoryRepo CategoryRepo) Service {
	return &service{
		categoryRepo: categoryRepo,
	}
}

func (svc *service) Create(name string) (*domain.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apperr.Validation("category name is required")
	}

	exists, err := svc.categoryRepo.ExistsByName(name, 0)
	if err != nil {
		return nil, apperr.WrapInternal("check category name", err)
	}
	if exists {
		return nil, apperr.Conflict("category name already exists")
	}

	cat, err := svc.categoryRepo.Create(domain.Category{Name: name})
	if err != nil {
		return nil, apperr.WrapInternal("create category", err)
	}

	return cat, nil
}

func (svc *service) Get(id int64) (*domain.Category, error) {
	cat, err := svc.categoryRepo.FindByID(id)
	if err != nil {
		return nil, apperr.WrapInternal("find category", err)
	}
	if cat == nil {
		return nil, apperr.NotFound("category not found")
	}
	return cat, nil
}

func (svc *service) Update(id int64, name string) (*domain.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apperr.Validation("category name is required")
	}

	existing, err := svc.categoryRepo.FindByID(id)
	if err != nil {
		return nil, apperr.WrapInternal("find category", err)
	}
	if existing == nil {
		return nil, apperr.NotFound("category not found")
	}

	exists, err := svc.categoryRepo.ExistsByName(name, id)
	if err != nil {
		return nil, apperr.WrapInternal("check category name", err)
	}
	if exists {
		return nil, apperr.Conflict("category name already exists")
	}

	cat, err := svc.categoryRepo.Update(domain.Category{ID: id, Name: name})
	if err != nil {
		return nil, apperr.WrapInternal("update category", err)
	}

	return cat, nil
}

func (svc *service) Delete(id int64) error {
	existing, err := svc.categoryRepo.FindByID(id)
	if err != nil {
		return apperr.WrapInternal("find category", err)
	}
	if existing == nil {
		return apperr.NotFound("category not found")
	}

	if err := svc.categoryRepo.Delete(id); err != nil {
		return apperr.WrapInternal("delete category", err)
	}

	return nil
}

func (svc *service) List(page, limit int64) ([]domain.Category, int64, error) {
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

	total, err := svc.categoryRepo.Count()
	if err != nil {
		return nil, 0, apperr.WrapInternal("count categories", err)
	}

	categories, err := svc.categoryRepo.List(limit, offset)
	if err != nil {
		return nil, 0, apperr.WrapInternal("list categories", err)
	}

	return categories, total, nil
}

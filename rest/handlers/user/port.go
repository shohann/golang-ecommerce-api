package user

import "github.com/shohann/golang-ecommerce-api/domain"

type Service interface {
	Create(user domain.User) (*domain.User, error)
	Login(email, pass string) (*LoginResponse, error)
	GetProfile(id int) (*domain.User, error)
	List(page, limit int64) ([]domain.User, int64, error)
}

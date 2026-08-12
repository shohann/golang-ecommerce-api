package user

import "github.com/shohann/golang-ecommerce-api/domain"

type Service interface {
	Create(user domain.User) (*domain.User, error)
	Login(email, pass string) (*domain.User, error)
}

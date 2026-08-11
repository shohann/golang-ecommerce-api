package user

import "github.com/shohann/golang-ecommerce-api/domain"

type UserRepo interface {
	Create(user domain.User) (*domain.User, error)
	FindAuthUser(email, pass string) (*domain.User, error)
}

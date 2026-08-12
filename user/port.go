package user

import "github.com/shohann/golang-ecommerce-api/domain"

type UserRepo interface {
	Create(user domain.User) (*domain.User, error)
	FindUserByEmail(email string) (*domain.User, error)
	CheckUniqueUser(email string) (bool, error)
}

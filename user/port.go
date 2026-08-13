package user

import "github.com/shohann/golang-ecommerce-api/domain"

type UserRepo interface {
	Create(user domain.User) (*domain.User, error)
	FindUserByEmail(email string) (*domain.User, error)
	CheckUniqueUser(email string) (bool, error)
	FindUserById(id int) (*domain.User, error)
	List(limit, offset int64) ([]domain.User, error)
	Count() (int64, error)
}

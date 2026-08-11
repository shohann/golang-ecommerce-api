package user

import (
	"github.com/shohann/golang-ecommerce-api/domain"
	userHandler "github.com/shohann/golang-ecommerce-api/rest/handlers/user"
)

type service struct {
	userRepo UserRepo
}

type Service interface {
	userHandler.Service
}

func NewService(userRepo UserRepo) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (svc *service) Create(user domain.User) (*domain.User, error) {
	usr, err := svc.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, nil
	}

	return usr, nil
}

func (svc *service) Find(email string, pass string) (*domain.User, error) {
	usr, err := svc.userRepo.FindAuthUser(email, pass)

	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, nil
	}

	return usr, nil
}

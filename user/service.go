package user

import (
	"github.com/shohann/golang-ecommerce-api/apperr"
	"github.com/shohann/golang-ecommerce-api/domain"
	userHandler "github.com/shohann/golang-ecommerce-api/rest/handlers/user"
	"golang.org/x/crypto/bcrypt"
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
	exist, err := svc.userRepo.CheckUniqueUser(user.Email)
	if err != nil {
		return nil, apperr.WrapInternal("check unique user", err)
	}

	if exist {
		return nil, apperr.Conflict("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperr.Internal("hash password", err)
	}

	user.Password = string(hashedPassword)

	usr, err := svc.userRepo.Create(user)
	if err != nil {
		return nil, apperr.WrapInternal("create user", err)
	}

	if usr == nil {
		return nil, nil
	}

	return usr, nil
}

func (svc *service) Find(email string, pass string) (*domain.User, error) {
	usr, err := svc.userRepo.FindAuthUser(email, pass)
	if err != nil {
		return nil, apperr.WrapInternal("find auth user", err)
	}

	return usr, nil
}

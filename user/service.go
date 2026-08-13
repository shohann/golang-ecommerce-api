package user

import (
	"fmt"

	"github.com/shohann/golang-ecommerce-api/apperr"
	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/domain"
	userHandler "github.com/shohann/golang-ecommerce-api/rest/handlers/user"
	"github.com/shohann/golang-ecommerce-api/util"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepo UserRepo
	cnf      *config.Config
}

type Service interface {
	userHandler.Service
}

func NewService(userRepo UserRepo, cnf *config.Config) Service {
	return &service{
		userRepo: userRepo,
		cnf:      cnf,
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

func (svc *service) Login(email string, pass string) (*userHandler.LoginResponse, error) {
	usr, err := svc.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, apperr.WrapInternal("find auth user", err)
	}

	if usr == nil {
		return nil, apperr.Conflict("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(pass))
	if err != nil {
		return nil, apperr.Conflict("invalid credentials")
	}

	accessToken, err := util.CreateJWT(svc.cnf.JWTSecretKey, util.Payload{
		Sub:      fmt.Sprint(usr.ID),
		FullName: usr.FullName,
		Email:    usr.Email,
		Role:     usr.Role,
	})
	if err != nil {
		return nil, apperr.Internal("create access token", err)
	}

	return &userHandler.LoginResponse{
		ID:          usr.ID,
		AccessToken: accessToken,
	}, nil
}

func (svc *service) GetProfile(id int) (*domain.User, error) {
	usr, err := svc.userRepo.FindUserById(id)
	if err != nil {
		return nil, apperr.WrapInternal("find user by id", err)
	}

	if usr == nil {
		return nil, apperr.NotFound("user not found")
	}

	return usr, nil
}

func (svc *service) List(page, limit int64) ([]domain.User, int64, error) {
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

	total, err := svc.userRepo.Count()
	if err != nil {
		return nil, 0, apperr.WrapInternal("count users", err)
	}

	users, err := svc.userRepo.List(limit, offset)
	if err != nil {
		return nil, 0, apperr.WrapInternal("list users", err)
	}

	return users, total, nil
}

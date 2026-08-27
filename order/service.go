package order

import (
	"github.com/shohann/golang-ecommerce-api/apperr"
	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/domain"
	cartRepo "github.com/shohann/golang-ecommerce-api/repo"
	orderHandler "github.com/shohann/golang-ecommerce-api/rest/handlers/order"
)

type service struct {
	cartItemRepo cartRepo.CartItemRepo
	cnf          *config.Config
}

type Service interface {
	orderHandler.Service
}

func NewService(
	cartItemRepo cartRepo.CartItemRepo,
	cnf *config.Config,
) Service {
	return &service{
		cartItemRepo: cartItemRepo,
		cnf:          cnf,
	}
}

func (svc *service) OrderCheckOut(userId int64) ([]domain.CartItem, error) {
	userCartItems, err := svc.cartItemRepo.GetCartByUserId(userId)

	if err != nil {
		return nil, apperr.WrapInternal("GetCartByUserId", err)
	}

	if userCartItems == nil {
		return nil, apperr.NotFound("cart item not found")
	}

	return userCartItems, nil
}

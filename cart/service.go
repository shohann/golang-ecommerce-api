package cart

import (
	"github.com/shohann/golang-ecommerce-api/apperr"
	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/domain"
	cartHandler "github.com/shohann/golang-ecommerce-api/rest/handlers/cart"
)

type service struct {
	cartItemRepo CartItemRepo
	cnf          *config.Config
}

type Service interface {
	cartHandler.Service
}

func NewService(cartItemRepo CartItemRepo, cnf *config.Config) Service {
	return &service{
		cartItemRepo: cartItemRepo,
		cnf:          cnf,
	}
}

func (svc *service) AddItem(cartItem domain.CartItem) (*domain.CartItem, error) {
	createdItem, err := svc.cartItemRepo.AddItem(cartItem)

	if err != nil {
		return nil, apperr.WrapInternal("create_product", err)
	}

	if createdItem == nil {
		return nil, nil
	}

	return createdItem, nil

}

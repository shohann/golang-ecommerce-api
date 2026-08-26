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

func (svc *service) AddItem(userID, productID int64, quantity int) (*domain.CartItem, error) {
	if quantity < 1 {
		return nil, apperr.Validation("quantity must be at least 1")
	}

	cartID, err := svc.cartItemRepo.GetOrCreateCartID(userID)
	if err != nil {
		return nil, apperr.WrapInternal("get or create cart", err)
	}

	createdItem, err := svc.cartItemRepo.AddItem(domain.CartItem{
		CartID:    cartID,
		ProductID: productID,
		Quantity:  quantity,
	})
	if err != nil {
		return nil, apperr.WrapInternal("add cart item", err)
	}

	if createdItem == nil {
		return nil, nil
	}

	return createdItem, nil
}

func (svc *service) GetCartByUserId(userId int64) ([]domain.CartItem, error) {
	userCartItems, err := svc.cartItemRepo.GetCartByUserId(userId)

	if err != nil {
		return nil, apperr.WrapInternal("GetCartByUserId", err)
	}

	if userCartItems == nil {
		return nil, apperr.NotFound("user not found")
	}

	return userCartItems, nil

}

func (svc *service) UpdateItem(userID, itemID int64, quantity int) (*domain.CartItem, error) {
	if quantity < 1 {
		return nil, apperr.Validation("quantity must be at least 1")
	}

	item, err := svc.cartItemRepo.UpdateQuantity(userID, itemID, quantity)
	if err != nil {
		return nil, apperr.WrapInternal("update cart item", err)
	}

	return item, nil
}

func (svc *service) DeleteItem(userID, itemID int64) error {
	if err := svc.cartItemRepo.DeleteItem(userID, itemID); err != nil {
		return apperr.WrapInternal("delete cart item", err)
	}

	return nil
}

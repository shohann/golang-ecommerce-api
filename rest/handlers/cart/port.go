package cart

import "github.com/shohann/golang-ecommerce-api/domain"

type Service interface {
	AddItem(userID, productID int64, quantity int) (*domain.CartItem, error)
	GetCartByUserId(userId int64) ([]domain.CartItem, error)
	UpdateItem(userID, itemID int64, quantity int) (*domain.CartItem, error)
	DeleteItem(userID, itemID int64) error
}

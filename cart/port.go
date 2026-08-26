package cart

import "github.com/shohann/golang-ecommerce-api/domain"

type CartItemRepo interface {
	AddItem(cartItem domain.CartItem) (*domain.CartItem, error)
	GetCartByUserId(userId int64) ([]domain.CartItem, error)
	UpdateQuantity(userID, itemID int64, quantity int) (*domain.CartItem, error)
	DeleteItem(userID, itemID int64) error
	GetOrCreateCartID(userID int64) (int64, error)
}

package cart

import "github.com/shohann/golang-ecommerce-api/domain"

type CartItemRepo interface {
	AddItem(cartItem domain.CartItem) (*domain.CartItem, error)
}

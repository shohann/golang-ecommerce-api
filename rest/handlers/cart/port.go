package cart

import "github.com/shohann/golang-ecommerce-api/domain"

type Service interface {
	AddItem(cartItem domain.CartItem) (*domain.CartItem, error)
}

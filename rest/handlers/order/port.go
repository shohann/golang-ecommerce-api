package order

import "github.com/shohann/golang-ecommerce-api/domain"

type Service interface {
	OrderCheckOut(userId int64) ([]domain.CartItem, error)
}
